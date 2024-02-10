package main

import (
	"log"
	"runtime"
	"unsafe"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
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

	window, err := glfw.CreateWindow(width, height, "Chaos", nil, nil)
	if err != nil {
		log.Fatalln("failed to create window:", err)
	}

	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		log.Fatalln("failed to initialize OpenGL", err)
	}

	return window
}

func scene_setup() {
	gl.Enable(gl.DEPTH_TEST)

	gl.ClearColor(0, 0, 0.0, 0.0)
	gl.ClearDepth(1)
	gl.DepthFunc(gl.LEQUAL)

	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	f := float64(width)/height - 1
	gl.Frustum(-1-f, 1+f, -1, 1, 1.0, 1000000.0)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
}

type Trail struct {
	Position *mgl64.Vec3
	Length   int
	Width    float32
	Color    mgl32.Vec3
	Curve    []mgl32.Vec3
}

func (t *Trail) Update(scale float64) {
	p := (*t.Position).Mul(scale)
	t.Curve = append(t.Curve, mgl32.Vec3{float32(p[0]), float32(p[1]), float32(p[2])})
	if len(t.Curve) > t.Length {
		t.Curve = t.Curve[1:]
	}
}

func (t *Trail) Draw() {
	gl.LineWidth(t.Width)
	gl.Begin(gl.LINE_STRIP)

	for i := 0; i < len(t.Curve); i++ {
		fade := float64(i) / float64(len(t.Curve))
		gl.Color4f(t.Color[0], t.Color[1], t.Color[2], float32(fade))
		gl.Vertex3f(t.Curve[i][0], t.Curve[i][1], t.Curve[i][2])
	}

	gl.End()
}

func main() {
	window := glfw_setup()
	defer glfw.Terminate()

	scene_setup()

	var colors []mgl32.Vec3
	var s Simulation
	s.Time = 10000.0
	s.Scale = 0.0000000005
	s.Points, colors = constructSystem("solar_system.toml")

	trails := make([]Trail, len(s.Points))
	for i := range trails {
		trails[i].Length = 200
		trails[i].Width = 3
		trails[i].Color = colors[i]
		trails[i].Position = &s.Points[i].Position
		trails[i].Update(s.Scale)
	}

	var c Controls
	c.Window = *window
	c.P.Orientation = mgl64.Vec3{0, 0, 0}
	c.P.Position = mgl64.Vec3{0, 0, 0}
	c.Inertia = mgl64.Vec3{0, 0, 0}
	c.Setup()

	for !window.ShouldClose() {
		// static behaviour
		s.Step()
		for i := range trails {
			trails[i].Update(s.Scale)
		}

		c.Handle(&s)
		m := c.P.Matrix()
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.MatrixMode(gl.MODELVIEW)
		gl.LoadMatrixd((*float64)(unsafe.Pointer(&m)))

		for i := 0; i < len(trails); i++ {
			trails[i].Draw()
		}

		draw_fadenkreuz(&c.P, 20.0)
		fillVoid()
		window.SwapBuffers()
	}
}

func fillVoid() {
	n := 20
	s := float32(10.0)
	gl.Begin(gl.QUADS)

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			for k := 0; k < n; k++ {
				gl.Color4f(float32(i)/float32(n), float32(j)/float32(n), float32(k)/float32(n), 1.0)
				gl.Vertex3f(float32(i)*s+float32(0), float32(j)*s+float32(0), float32(k)*s)
				gl.Vertex3f(float32(i)*s+float32(0), float32(j)*s+float32(1), float32(k)*s)
				gl.Vertex3f(float32(i)*s+float32(1), float32(j)*s+float32(1), float32(k)*s)
				gl.Vertex3f(float32(i)*s+float32(1), float32(j)*s+float32(0), float32(k)*s)
			}
		}
	}
	gl.End()
}

func draw_fadenkreuz(p *Pov, d float64) {
	gl.Begin(gl.LINES)

	var c [4]float32
	gl.Color4f(1.0, 0.0, 0.0, 1.0)

	p.FreeMove(mgl64.Vec3{0, 0, -2 * d})
	base := p.Position
	p.FreeMove(mgl64.Vec3{0, 0, 2 * d})

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
