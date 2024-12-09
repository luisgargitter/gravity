package main

import (
	"fmt"
	"github.com/luisgargitter/numerics"
	_ "image/jpeg"
	"log"
	"math"
	"runtime"
	"strings"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl64"
)

const MouseSensi = 0.0005
const (
	RenderHeight = 1200
	RenderWidth  = 1600
)
const (
	PlotHeight = 400
	PlotWidth  = 600
)

const glCorrectionScale = 10e-9

const fpsTarget = 60.0

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	fmt.Println("Initialization...")
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	renderer := InitRenderer(RenderWidth, RenderHeight)
	plotter := InitPlotter(PlotWidth, PlotHeight, 200)
	defer glfw.Terminate()

	renderer.win.MakeContextCurrent()

	p := Pov{mgl64.Vec3{0, 0, 20e9}, mgl64.Vec3{0, 0, 1}, mgl64.Vec3{0, 1, 0}}

	var c Controls
	c.Window = *renderer.win
	c.P = p
	c.Velocity = mgl64.Vec3{0, 0, 0}
	c.Acceleration = 1000
	c.Resistance = 1.0
	c.PlanetIndex = 3
	c.Setup()

	sphere_vao := loadSphere(5, 1.0)

	fmt.Println("Loading Planetary System...")
	particles, radii, textures, names := constructSystem("solar_system.toml")
	fmt.Println("Planetary System Loaded.")

	objects := make([]Object, len(particles))
	for i := range particles {
		pos := particles[i].Position.Mul(glCorrectionScale)
		r := radii[i] * 10 * glCorrectionScale
		t := mgl64.Translate3D(pos[0], pos[1], pos[2]).Mul4(mgl64.Scale3D(r, r, r)).Mul4(mgl64.HomogRotate3D(-math.Pi/2, mgl64.Vec3{1, 0, 0}))
		objects[i] = Object{t, textures[i], sphere_vao}
	}
	sVerts := []Particle{particles[3], particles[9]}
	sLinks := []Edge[Link]{Edge[Link]{0, 1, Link{sVerts[0].Position.Sub(sVerts[1].Position).Len() / 2, 100, 0}}}
	satellite := Softbody{sVerts, sLinks}

	timeScale := 100000.0

	particlesRK4W := numerics.RK4Workspace[ParticleSystem]{
		Add: particleSystemAdd,
		Mul: particleSystemMul,
		D:   make(ParticleSystem, len(particles)),
		K1:  make(ParticleSystem, len(particles)),
		K2:  make(ParticleSystem, len(particles)),
		K3:  make(ParticleSystem, len(particles)),
		K4:  make(ParticleSystem, len(particles)),
	}

	satelliteRK4W := numerics.RK4Workspace[Softbody]{
		Add: softbodyAdd,
		Mul: softbodyMul,
		D:   Softbody{make([]Particle, len(satellite.vertices)), make([]Edge[Link], len(satellite.edges))},
		K1:  Softbody{make([]Particle, len(satellite.vertices)), make([]Edge[Link], len(satellite.edges))},
		K2:  Softbody{make([]Particle, len(satellite.vertices)), make([]Edge[Link], len(satellite.edges))},
		K3:  Softbody{make([]Particle, len(satellite.vertices)), make([]Edge[Link], len(satellite.edges))},
		K4:  Softbody{make([]Particle, len(satellite.vertices)), make([]Edge[Link], len(satellite.edges))},
	}

	camera := Camera{
		&c.P.Position, &c.P.Orientation, &c.P.Up,
		math.Pi / 4.0, float64(RenderWidth) / float64(RenderHeight),
		1e7,
		1.0e12,
	}
	scene := Scene{&camera, objects}

	var cpuTime, gpuTime, deltaTime float64

	info := Info{
		&c.P.Position, &c.Velocity, &c.P.Orientation,
		&cpuTime, &gpuTime, &deltaTime,
		&names, &c.Locked, &c.PlanetIndex,
	}

	i := 0
	for ; !renderer.win.ShouldClose() && !plotter.win.ShouldClose(); i++ {
		t := glfw.GetTime()

		if i%fpsTarget == 0 {
			i = 0
			info.Print()
		}

		if i%(fpsTarget/6) == 0 {
			plotter.Update()
		}

		// static behaviour
		numerics.RK4(&particlesRK4W, dParticleSystem, deltaTime*timeScale, &particles, &particles)
		numerics.RK4(&satelliteRK4W, dSoftbody, deltaTime*timeScale, &satellite, &satellite)

		c.Handle(particles, deltaTime*timeScale)

		gl.ClearColor(0, 0, 0, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		for i := range scene.objects {
			pos := particles[i].Position.Mul(glCorrectionScale)
			r := radii[i] * 10 * glCorrectionScale
			scene.objects[i].Transform = mgl64.Translate3D(pos[0], pos[1], pos[2]).Mul4(mgl64.Scale3D(r, r, r).Mul4(mgl64.HomogRotate3D(-math.Pi/2, mgl64.Vec3{1, 0, 0})))
		}

		cpuTime = glfw.GetTime() - t

		scene.Draw(renderer.viewUni)
		renderer.win.SwapBuffers()

		plotter.Draw()
		renderer.win.MakeContextCurrent()

		gpuTime = glfw.GetTime() - (t + cpuTime)

		sleepTime := 1.0/fpsTarget - (cpuTime + gpuTime)

		time.Sleep(time.Duration(1000.0*sleepTime) * time.Millisecond)
		deltaTime = glfw.GetTime() - t
	}
	fmt.Print("\n")
}

func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

var vertexShaderSource = /*`
#version 330 core
layout (location = 0) in vec3 vert;

void main()
{
    gl_Position = vec4(vert, 1.0);
}
` + "\x00"
*/

`
#version 330 core

uniform mat4 view;

in vec3 vert;
in vec2 uv;

out vec2 texcoord;

void main() {
    gl_Position = mat4(view) * vec4(vert, 1.0f);
	texcoord = uv;
}
` + "\x00"

var fragmentShaderSource = `
#version 330 core
out vec4 outputColor;
in vec2 texcoord;

uniform sampler2D tex;

void main()
{
	//outputColor = vec4(vec3(1/gl_FragCoord.z), 1.0);
	outputColor = texture(tex, texcoord);
    //outputColor = vec4(texcoord[0], texcoord[1], 0.0f, 1.0f);
}
` + "\x00"

/*
`
#version 330

uniform sampler2D tex;

in vec2 fragTexCoord;

out vec4 outputColor;

void main() {
    outputColor = texture(tex, fragTexCoord);
}
` + "\x00"
*/
