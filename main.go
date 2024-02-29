package main

import (
	"fmt"
	"log"
	"math"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl64"
)

const mouse_sensi = 0.005
const width, height = 1600, 1200

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

	if err := gl.Init(); err != nil {
		log.Fatalln("failed to initialize OpenGL", err)
	}

	return window
}

func gl_setup() {
	gl.Enable(gl.CULL_FACE)
	//gl.CullFace(gl.FRONT)

	//gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE) // wireframe
	gl.PolygonMode(gl.FRONT, gl.FILL)
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
	window := glfw_setup()
	defer glfw.Terminate()
	gl_setup()

	p := Pov{mgl64.Vec3{}, mgl64.Vec3{}}

	var c Controls
	c.Window = *window
	c.P = p
	c.Inertia = mgl64.Vec3{0, 0, 0}
	c.Acceleration = 10000000
	c.Resistance = 0.95
	c.Setup()

	var os []Object

	sphere_vao := loadSphere()

	var radii []float64
	var s Simulation
	s.Time = 10000.0
	s.Points, _, radii = constructSystem("solar_system.toml")
	for i := range s.Points {
		pos := s.Points[i].Position
		r := radii[i]
		os = append(os, Object{mgl64.Translate3D(pos[0], pos[1], pos[2]).Mul4(mgl64.Scale3D(r, r, r)), sphere_vao})
	}

	program, err := newProgram(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}
	gl.UseProgram(program)

	viewU := gl.GetUniformLocation(program, gl.Str("view\x00"))

	projection := mgl64.Perspective(math.Pi/2, float64(width)/float64(height), 0.1, 1.0e12)
	camera := Camera{projection, &c.P}
	scene := Scene{&camera, os}

	var cpuEnd, cpuStart float64
	var gpuEnd, gpuStart float64

	var info Info
	info.Position = &c.P.Position
	info.Orientation = &c.P.Orientation
	info.CpuEnd = &cpuEnd
	info.CpuStart = &cpuStart
	info.GpuEnd = &gpuEnd
	info.GpuStart = &gpuStart

	for !window.ShouldClose() {
		cpuStart = glfw.GetTime()
		// static behaviour
		s.Step()
		c.Handle(&s)

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		for i := range scene.Os {
			pos := s.Points[i].Position
			r := radii[i]
			scene.Os[i].Transform = mgl64.Translate3D(pos[0], pos[1], pos[2]).Mul4(mgl64.Scale3D(r, r, r))
		}

		cpuEnd = glfw.GetTime()
		gpuStart = cpuEnd

		scene.Draw(viewU)
		window.SwapBuffers()

		gpuEnd = glfw.GetTime()

		info.Print()
	}
}

func draw_fadenkreuz(p *Pov, d float64) {
	gl.LineWidth(1)
	gl.Begin(gl.LINES)

	var c [4]float32
	gl.Color4f(1.0, 0.0, 0.0, 1.0)

	p.FreeMove(mgl64.Vec3{0, 0, 5 * d})
	base := p.Position
	p.FreeMove(mgl64.Vec3{0, 0, -5 * d})

	for i := 0; i < 3; i++ {
		c = [4]float32{0.0, 0.0, 0.0, 1.0}
		c[i] = 1.0
		gl.Color4f(c[0], c[1], c[2], c[3])
		t := base
		gl.Vertex3f(float32(t[0]), float32(t[1]), float32(t[2]))
		t[i] += d
		gl.Vertex3f(float32(t[0]), float32(t[1]), float32(t[2]))
	}
	gl.End()
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

void main() {
    gl_Position = mat4(view) * vec4(vert, 1);
}
` + "\x00"

var fragmentShader = `
#version 410 core
out vec4 outputColor;

void main()
{
    outputColor = vec4(1.0f, 0.5f, 0.2f, 1.0f);
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
