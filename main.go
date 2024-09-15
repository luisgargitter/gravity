package main

import (
	"fmt"
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

const mouse_sensi = 0.0005
const width, height = 800, 600

const glCorrectionScale = 10e-9

const fpsTarget = 60

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func glfw_setup() *glfw.Window {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	window, err := glfw.CreateWindow(width, height, "Gravity", nil, nil)
	if err != nil {
		log.Fatalln("failed to create window:", err)
	}

	window.MakeContextCurrent()
	glfw.SwapInterval(0) // for controlling framerate

	if err := gl.Init(); err != nil {
		log.Fatalln("failed to initialize OpenGL", err)
	}

	return window
}

func gl_setup() {
	gl.Viewport(0, 0, width, height)

	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.FRONT)

	//gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE) // wireframe
	gl.PolygonMode(gl.BACK, gl.FILL)
}

func loadSphere() VAO {
	cube := Cube()
	for i := 0; i < 5; i++ {
		cube.Enhance()
	}
	cube.PuffUp(1)

	return cube.Load()
}

func main() {
	fmt.Println("Initialization...")
	window := glfw_setup()
	defer glfw.Terminate()
	gl_setup()

	p := Pov{mgl64.Vec3{0, 0, 20e9}, mgl64.Vec3{}}

	var c Controls
	c.Window = *window
	c.P = p
	c.Inertia = mgl64.Vec3{0, 0, 0}
	c.Acceleration = 100000
	c.Resistance = 1.0
	c.PlanetIndex = 3
	c.Setup()

	var objects []Object

	sphere_vao := loadSphere()

	var radii []float64
	var textures []uint32
	var names []string
	var s Simulation
	s.Time = 100.0
	fmt.Println("Loading Planetary System...")
	s.Points, radii, textures, names = constructSystem("solar_system.toml")
	fmt.Println("Planetary System Loaded.")

	for i := range s.Points {
		pos := s.Points[i].Position.Mul(glCorrectionScale)
		r := radii[i] * 10 * glCorrectionScale
		objects = append(objects,
			Object{mgl64.Translate3D(pos[0], pos[1], pos[2]).Mul4(mgl64.Scale3D(r, r, r)).Mul4(mgl64.HomogRotate3D(-math.Pi/2, mgl64.Vec3{1, 0, 0})),
				textures[i],
				sphere_vao},
		)
	}

	fmt.Println("Compiling Shaders...")
	program, err := newProgram(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}
	fmt.Println("Compilation Done.")
	gl.UseProgram(program)

	viewU := gl.GetUniformLocation(program, gl.Str("view\x00"))

	projection := mgl64.Perspective(math.Pi/4, float64(width)/float64(height), 0.1, 1.0e12*glCorrectionScale)
	camera := Camera{projection, &c.P}
	scene := Scene{&camera, objects}

	var cpuTime, gpuTime, deltaTime float64

	var info Info
	info.Position = &c.P.Position
	info.Inertia = &c.Inertia
	info.Orientation = &c.P.Orientation
	info.CpuTime = &cpuTime
	info.GpuTime = &gpuTime
	info.DeltaTime = &deltaTime
	info.Planets = &names
	info.Locked = &c.Locked
	info.PlanetIndex = &c.PlanetIndex

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	i := 0
	for ;!window.ShouldClose(); i++ {
		if i % fpsTarget == 0 {
			i = 0
			info.Print()
		}

		deltaTime = glfw.GetTime()
		cpuTime = deltaTime
		// static behaviour
			s.Step()
		
		c.Handle(&s)

		gl.ClearColor(0, 0, 0, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		for i := range scene.Os {
			pos := s.Points[i].Position.Mul(glCorrectionScale)
			r := radii[i] * 10 * glCorrectionScale
			scene.Os[i].Transform = mgl64.Translate3D(pos[0], pos[1], pos[2]).Mul4(mgl64.Scale3D(r, r, r).Mul4(mgl64.HomogRotate3D(-math.Pi/2, mgl64.Vec3{1, 0, 0})))
		}

		gpuTime = glfw.GetTime()
		cpuTime = gpuTime - cpuTime

		scene.Draw(viewU)
		window.SwapBuffers()
		gpuTime = glfw.GetTime() - gpuTime

		sleepTime := 1.0/fpsTarget - (cpuTime + gpuTime)
		time.Sleep(time.Duration(sleepTime) * time.Second)
		deltaTime = glfw.GetTime() - deltaTime
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

var vertexShader = /*`
#version 330 core
layout (location = 0) in vec3 vert;

void main()
{
    gl_Position = vec4(vert, 1.0);
}
` + "\x00"
*/

`
#version 410

uniform dmat4 view;

in vec3 vert;
in vec2 uv;

out vec2 texcoord;

void main() {
    gl_Position = mat4(view) * vec4(vert, 1.0f);
	texcoord = uv;
}
` + "\x00"

var fragmentShader = `
#version 410 core
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
