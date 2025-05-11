package graphics

import (
	"fmt"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	_ "image/jpeg"
	"log"
	"os"
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v2.1/gl"
)

func glSetup(win *glfw.Window) uint32 {
	width, height := win.GetSize()
	gl.Viewport(0, 0, int32(width), int32(height))

	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.FRONT)

	//gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINES) // wireframe
	gl.PolygonMode(gl.BACK, gl.FILL)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	fmt.Println("Compiling Shaders...")
	program, err := newProgram("graphics/shaders/vert.vert", "graphics/shaders/distance.frag")
	if err != nil {
		panic(err)
	}
	fmt.Println("Compilation Done.")
	gl.UseProgram(program)

	return program
}

func BindRenderer(win *glfw.Window) uint32 {
	win.MakeContextCurrent()
	glfw.SwapInterval(1) // vsync (set to zero for unlimited framerate

	if err := gl.Init(); err != nil {
		log.Fatalln("failed to initialize OpenGL", err)
	}

	return glSetup(win)
}

type VBO uint32
type EBO struct {
	ebo   uint32
	count int32
}
type VAO struct {
	vao   uint32
	count int32
}

type Object struct {
	Mesh        *Mesh
	Position    mgl32.Vec3
	Orientation mgl32.Quat
}

func constructVBO(vertices []mgl32.Vec3) VBO {
	var r uint32
	a := make([][3]float32, len(vertices))
	for i := range vertices {
		v := vertices[i][:]
		a[i] = [3]float32{v[0], v[1], v[2]}
	}

	gl.GenBuffers(1, &r)
	gl.BindBuffer(gl.ARRAY_BUFFER, r)
	gl.BufferData(
		gl.ARRAY_BUFFER,
		int(unsafe.Sizeof(a[0]))*len(a),
		unsafe.Pointer(&a[0]),
		gl.STATIC_DRAW,
	)

	return VBO(r)
}

func constructEBO(faces []Surface) EBO {
	var r uint32
	gl.GenBuffers(1, &r)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, r)
	gl.BufferData(
		gl.ELEMENT_ARRAY_BUFFER,
		int(unsafe.Sizeof(faces[0]))*len(faces),
		unsafe.Pointer(&faces[0]),
		gl.STATIC_DRAW,
	)

	return EBO{r, int32(len(faces))}
}

func constructVAO(vbo VBO, ebo EBO) VAO {
	var r uint32
	gl.GenVertexArrays(1, &r)
	gl.BindVertexArray(r)

	gl.BindBuffer(gl.ARRAY_BUFFER, uint32(vbo))
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo.ebo)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, int32(unsafe.Sizeof([3]float32{})), nil)
	gl.EnableVertexAttribArray(0)

	gl.BindVertexArray(0) // unbind

	return VAO{r, ebo.count}
}

func (v *VAO) Draw() {
	gl.BindVertexArray(v.vao)
	gl.DrawElements(gl.TRIANGLES, int32(v.count*3), gl.UNSIGNED_INT, nil)
	gl.BindVertexArray(0)
}

func newProgram(vertexShaderFile, fragmentShaderFile string) (uint32, error) {
	vertexShaderSource, err := os.ReadFile(vertexShaderFile)
	if err != nil {
		return 0, err
	}
	fragmentShaderSource, err := os.ReadFile(fragmentShaderFile)
	if err != nil {
		return 0, err
	}

	vertexShader, err := compileShader(string(vertexShaderSource)+"\x00", gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(string(fragmentShaderSource)+"\x00", gl.FRAGMENT_SHADER)
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
