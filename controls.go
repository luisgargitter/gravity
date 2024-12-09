package main

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl64"
	"log"
)

type Pov struct {
	Position    mgl64.Vec3
	Orientation mgl64.Vec3
	Up          mgl64.Vec3
}

func (p *Pov) FPSLook(delta mgl64.Vec2) {
	qx := mgl64.QuatRotate(delta[0], p.Up)
	qy := mgl64.QuatRotate(delta[1], p.Orientation.Cross(p.Up))
	p.Orientation = qx.Mul(qy).Rotate(p.Orientation)
}

func (p *Pov) FreeMove(delta mgl64.Vec3) {
	oz := p.Orientation
	ox := p.Orientation.Cross(p.Up)
	oy := ox.Cross(p.Orientation)
	om := mgl64.Mat3FromCols(ox, oy, oz)
	p.Position = p.Position.Add(om.Mul3x1(delta)) // for now
}

func (p *Pov) FPSMove(delta mgl64.Vec3) {
	ox := p.Orientation.Cross(p.Up)
	oy := p.Up
	oz := oy.Cross(ox)
	om := mgl64.Mat3FromCols(ox, oy, oz)
	p.Position = p.Position.Add(om.Mul3x1(delta)) // for now
}

type Controls struct {
	Window       glfw.Window
	Mouse        mgl64.Vec2
	P            Pov
	Velocity     mgl64.Vec3
	Acceleration float64
	Resistance   float64
	Locked       bool
	PlanetIndex  int
}

func (c *Controls) Setup() {
	c.Window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	if !glfw.RawMouseMotionSupported() {
		log.Fatalln("raw mouse motion not supported")
	}
	c.Window.SetInputMode(glfw.RawMouseMotion, glfw.True)

	c.Mouse[0], c.Mouse[1] = c.Window.GetCursorPos()
}

func (c *Controls) Handle(particles []Particle, dt float64) {
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
	stop := c.Window.GetKey(glfw.KeyL)

	dV := c.Acceleration
	if sw == glfw.Press {
		c.Velocity[2] += dV
	}
	if ss == glfw.Press {
		c.Velocity[2] -= dV
	}
	if sd == glfw.Press {
		c.Velocity[0] += dV
	}
	if sa == glfw.Press {
		c.Velocity[0] -= dV
	}
	if up == glfw.Press {
		c.Velocity[1] += dV
	}
	if down == glfw.Press {
		c.Velocity[1] -= dV
	}
	if q == glfw.Press {
		c.Window.SetShouldClose(true)
	}
	if stop == glfw.Press {
		c.Velocity = mgl64.Vec3{0, 0, 0}
	}

	if lock == glfw.Press {
		c.Locked = true
		planet := particles[c.PlanetIndex]

		c.P.FreeMove(c.Velocity.Mul(dt))
		c.P.Position = c.P.Position.Add(planet.Velocity.Mul(dt))

		c.P.Orientation = planet.Position.Sub(c.P.Position).Normalize()

	} else {
		c.P.FPSMove(c.Velocity.Mul(dt))
		c.P.FPSLook(c.Mouse.Sub(mouse).Mul(MouseSensi))
	}
	if lock == glfw.Release && c.Locked {
		c.Locked = false
		c.PlanetIndex = (c.PlanetIndex + 1) % len(particles)
	}

	c.Velocity = c.Velocity.Mul(c.Resistance) // for smooth movement (Kondensator-Ladekurve)
	c.Mouse = mouse
}
