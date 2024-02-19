package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
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

func scene_setup() {
	gl.Enable(gl.DEPTH_TEST)

	gl.ClearColor(0, 0, 0.0, 0.0)
	gl.ClearDepth(1)
	gl.DepthFunc(gl.LEQUAL)

	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	f := float64(width)/height - 1
	gl.Frustum(-1-f, 1+f, -1, 1, 1.0, 1.0e12)
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
		fade := float64(i) / float64(t.Length)
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
	var radii []float64
	var s Simulation
	s.Time = 10000.0
	s.Scale = 0.000000005
	s.Points, colors, radii = constructSystem("solar_system.toml")

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
	c.P.Orientation = mgl64.Vec3{0, math.Pi, 0}
	c.P.Position = mgl64.Vec3{0, 0, 0}
	c.Inertia = mgl64.Vec3{0, 0, 0}
	c.Acceleration = 1000000.0
	c.Resistance = 0.95
	c.Setup()

	cube := Cube()
	cube.Colors = []mgl32.Vec3{
		{rand.Float32(), rand.Float32(), rand.Float32()},
		{rand.Float32(), rand.Float32(), rand.Float32()},
		{rand.Float32(), rand.Float32(), rand.Float32()},
		{rand.Float32(), rand.Float32(), rand.Float32()},
		{rand.Float32(), rand.Float32(), rand.Float32()},
		{rand.Float32(), rand.Float32(), rand.Float32()},
		{rand.Float32(), rand.Float32(), rand.Float32()},
		{rand.Float32(), rand.Float32(), rand.Float32()},
	}

	for i := 0; i < 4; i++ {
		cube.Enhance()
	}
	cube.PuffUp(1)

	var sphereBuffer uint32
	gl.GenBuffers(1, &sphereBuffer)

	avgStarDis := 9.461e+18

	var stars []mgl32.Vec3
	for i := 0; i < 1000; i++ {
		stars = append(stars, mgl32.SphericalToCartesian(
			rand.Float32()*float32(avgStarDis*s.Scale),
			float32(math.Asin(2*rand.Float64()-1)+math.Pi/2),
			rand.Float32()*float32(2*math.Pi),
		))
	}

	for !window.ShouldClose() {
		cpuStart := glfw.GetTime()
		// static behaviour
		s.Step()
		for i := range trails {
			trails[i].Update(s.Scale)
		}

		c.Handle(&s)
		cpuEnd := glfw.GetTime()
		gpuStart := cpuEnd

		m := c.P.Matrix()
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.MatrixMode(gl.MODELVIEW)
		gl.LoadMatrixd((*float64)(unsafe.Pointer(&m)))
		//gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE) // wireframe
		gl.PolygonMode(gl.FRONT, gl.FILL)

		for i := 0; i < len(trails); i++ {
			t := s.Points[i].Position.Mul(s.Scale)
			r := radii[i] * s.Scale
			gl.Translated(t[0], t[1], t[2])
			gl.Scaled(r, r, r)
			tc := trails[i].Color
			for j := range cube.Colors {
				cube.Colors[j] = tc
			}
			cube.Draw()
			r = 1 / r
			gl.Scaled(r, r, r)
			gl.Translated(-t[0], -t[1], -t[2])

			trails[i].Draw()
		}

		gl.Begin(gl.POINTS)
		gl.Color4f(1, 1, 1, 1)
		for _, s := range stars {
			gl.Vertex3f(s[0], s[1], s[2])
		}
		gl.End()

		draw_fadenkreuz(&c.P, 1.0)
		// fillVoid()
		window.SwapBuffers()
		gpuEnd := glfw.GetTime()
		cpuTime := (cpuEnd - cpuStart)
		gpuTime := (gpuEnd - gpuStart)
		fps := 1 / (cpuTime + gpuTime)
		fmt.Printf("\rPosition: (%.2f, %.2f, %.2f) Orientation: (%.2f, %.2f, %.2f) CPU: %.2f ms, GPU: %.2f ms, FPS: %.2f ",
			c.P.Position[0], c.P.Position[1], c.P.Position[2],
			c.P.Orientation[0], c.P.Orientation[1], c.P.Orientation[2],
			cpuTime*1000, gpuTime*1000, fps)
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
