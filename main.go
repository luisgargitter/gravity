package main

import (
	"fmt"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"gravity/graphics"
	"gravity/physics"
	_ "image/jpeg"
	"log"
	"math"
	"runtime"
	"time"
)

const MouseSensi = 0.0005
const (
	RenderHeight = 600
	RenderWidth  = 800
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
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	win, err := glfw.CreateWindow(RenderWidth, RenderHeight, "Gravity", nil, nil)
	if err != nil {
		log.Fatalln("Failed to create window:", err)
	}
	program := graphics.BindRenderer(win)
	plotter := graphics.InitPlotter(PlotWidth, PlotHeight, 200)
	defer glfw.Terminate()

	win.MakeContextCurrent()

	p := Pov{mgl64.Vec3{0, 0, 20e9}, mgl64.Vec3{0, 0, 1}, mgl64.Vec3{0, 1, 0}}

	var c Controls
	c.Window = *win
	c.P = p
	c.Velocity = mgl64.Vec3{0, 0, 0}
	c.Acceleration = 1000
	c.Resistance = 1.0
	c.PlanetIndex = 3
	c.Setup()

	fmt.Println("Loading Planetary System...")
	particles, names := constructSystem("solar_system.toml")
	fmt.Println("Planetary System Loaded.")

	sphere := graphics.Sphere(5)

	objects := make([]graphics.Object, len(particles))
	for i := range particles {
		pos64 := particles[i].Position.Mul(glCorrectionScale)
		pos := mgl32.Vec3{float32(pos64[0]), float32(pos64[1]), float32(pos64[2])}
		objects[i] = graphics.Object{&sphere, pos, mgl32.QuatIdent()}
	}

	timeScale := 100000.0

	particlesRK4W := physics.ParticleSystemRK4W(particles)

	cpos := c.P.Position.Mul(glCorrectionScale)
	cpos32 := mgl32.Vec3{float32(cpos[0]), float32(cpos[1]), float32(cpos[2])}
	cor := c.P.Orientation
	cor32 := mgl32.Vec3{float32(cor[0]), float32(cor[1]), float32(cor[2])}
	cup := c.P.Up
	cup32 := mgl32.Vec3{float32(cup[0]), float32(cup[1]), float32(cup[2])}
	camera := graphics.Camera{
		cpos32, cor32, cup32,
		math.Pi / 4.0, float32(RenderWidth) / float32(RenderHeight),
		1e7,
		1.0e12,
	}
	scene := graphics.SceneInit(camera, program, objects)

	var cpuTime, gpuTime, deltaTime float64

	info := Info{
		&c.P.Position, &c.Velocity, &c.P.Orientation,
		&cpuTime, &gpuTime, &deltaTime,
		&names, &c.Locked, &c.PlanetIndex,
	}

	earthLastPos := particles[3].Position
	earthDeltaPos := mgl64.Vec3{}
	plotDiff := mgl64.Vec3{}

	plotter.Attach(&plotDiff[0])
	plotter.Attach(&plotDiff[2])

	i := 0
	for ; !win.ShouldClose() && !plotter.Win.ShouldClose(); i++ {
		t := glfw.GetTime()

		if i%fpsTarget == 0 {
			i = 0
			info.Print()
		}

		if i%(fpsTarget/6) == 0 {
			plotter.Update()
		}

		physics.RK4(particlesRK4W, deltaTime*timeScale)

		c.Handle(particles, deltaTime*timeScale)

		for i := range particles {
			pos64 := particles[i].Position.Mul(glCorrectionScale)
			pos := mgl32.Vec3{float32(pos64[0]), float32(pos64[1]), float32(pos64[2])}
			objects[i].Position = pos
		}

		cpuTime = glfw.GetTime() - t

		scene.Draw()
		win.SwapBuffers()

		earthDeltaPos = particles[3].Position.Sub(earthLastPos)
		earthLastPos = particles[3].Position
		plotDiff = earthDeltaPos.Mul(1 / (deltaTime * timeScale)).Sub(particles[3].Velocity)
		plotter.Draw()
		win.MakeContextCurrent()

		gpuTime = glfw.GetTime() - (t + cpuTime)

		sleepTime := 1.0/fpsTarget - (cpuTime + gpuTime)

		time.Sleep(time.Duration(1000.0*sleepTime) * time.Millisecond)
		deltaTime = glfw.GetTime() - t
	}
	fmt.Print("\n")
}
