package main

import (
	"fmt"
	"log"
	"math"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl64"
)

type Pov struct {
	Position    mgl64.Vec3
	Orientation mgl64.Vec3
}

func (p *Pov) FPSLook(delta mgl64.Vec2) {
	p.Orientation[1] += delta[0]
	p.Orientation[0] += delta[1]

	p.Orientation[1] = math.Mod(p.Orientation[1], 2*math.Pi)
	p.Orientation[0] = mgl64.Clamp(p.Orientation[0], -math.Pi/2, math.Pi/2)
}

func (p *Pov) FreeMove(delta mgl64.Vec3) {
	delta[2] = -delta[2] // OpenGL uses a right hand coordinate system.
	// To associate +z with forward movement we need to invert it.
	q := mgl64.AnglesToQuat(p.Orientation[2], p.Orientation[1], p.Orientation[0], mgl64.ZYX)
	p.Position = p.Position.Add(q.Rotate(delta))
}

func (p *Pov) Matrix() mgl64.Mat4 {
	po := mgl64.Translate3D(p.Position[0], p.Position[1], p.Position[2])
	or := mgl64.AnglesToQuat(p.Orientation[2], p.Orientation[1], p.Orientation[0], mgl64.ZYX).Mat4()

	return po.Mul4(or).Inv()
}

func (p *Pov) FPSMove(delta mgl64.Vec3) {
	t := p.Orientation
	p.Orientation = mgl64.Vec3{0, p.Orientation[1], 0}
	p.FreeMove(delta)
	p.Orientation = t
}

type Controls struct {
	Window  glfw.Window
	Mouse   mgl64.Vec2
	P       Pov
	Inertia mgl64.Vec3
}

func (c *Controls) Setup() {
	c.Window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	if glfw.RawMouseMotionSupported() == false {
		log.Fatalln("raw mouse motion not supported")
	}
	c.Window.SetInputMode(glfw.RawMouseMotion, glfw.True)

	c.Mouse[0], c.Mouse[1] = c.Window.GetCursorPos()
}

func (c *Controls) Handle(s *Simulation) {
	var mouse mgl64.Vec2
	glfw.PollEvents()
	// input dependent behaviour
	mouse[0], mouse[1] = c.Window.GetCursorPos()

	sw := c.Window.GetKey(glfw.KeyW)
	sa := c.Window.GetKey(glfw.KeyA)
	ss := c.Window.GetKey(glfw.KeyS)
	sd := c.Window.GetKey(glfw.KeyD)
	up := c.Window.GetKey(glfw.KeySpace)
	down := c.Window.GetKey(glfw.KeyLeftShift)
	q := c.Window.GetKey(glfw.KeyQ)
	lock := c.Window.GetKey(glfw.KeyTab)

	if sw == glfw.Press {
		c.Inertia[2] += 0.01
	}
	if ss == glfw.Press {
		c.Inertia[2] -= 0.01
	}
	if sd == glfw.Press {
		c.Inertia[0] += 0.01
	}
	if sa == glfw.Press {
		c.Inertia[0] -= 0.01
	}
	if up == glfw.Press {
		c.Inertia[1] += 0.01
	}
	if down == glfw.Press {
		c.Inertia[1] -= 0.01
	}
	if q == glfw.Press {
		c.Window.SetShouldClose(true)
	}

	if lock == glfw.Press {
		c.P.FreeMove(c.Inertia)
		c.P.Position = c.P.Position.Add(s.Points[3].Inertia.Mul(s.Time * s.Scale))

		t := c.P.Position.Sub(s.Points[3].Position.Mul(s.Scale))
		_, theta, phi := mgl64.CartesianToSpherical(mgl64.Vec3{t[0], t[2], t[1]})

		c.P.Orientation = mgl64.Vec3{theta - math.Pi/2, -phi + math.Pi/2, 0}
	} else {
		c.P.FPSMove(c.Inertia)
		c.P.FPSLook(c.Mouse.Sub(mouse).Mul(mouse_sensi))
	}

	c.Inertia = c.Inertia.Mul(0.99) // for smooth movement (Kondensator-Ladekurve)
	c.Mouse = mouse

	fmt.Printf("\rPosition: (%.2f, %.2f, %.2f) Orientation: (%.2f, %.2f, %.2f)",
		c.P.Position[0], c.P.Position[1], c.P.Position[2],
		c.P.Orientation[0], c.P.Orientation[1], c.P.Orientation[2])
}
