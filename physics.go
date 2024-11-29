package main

import (
	"github.com/go-gl/mathgl/mgl64"
	"math"
)

const (
	G    = 6.6743015e-11
	Eps0 = 8.8541878128e-12
)

type Particle struct {
	Position mgl64.Vec3
	Velocity mgl64.Vec3
	Mass     float64
	Charge   float64
}

func (p *Particle) GravitationalForceV(a *Particle) mgl64.Vec3 {
	deltaPosition := a.Position.Sub(p.Position)
	distanceSquared := deltaPosition.LenSqr()
	Fg := G * p.Mass * a.Mass / distanceSquared
	Fc := p.Charge * a.Charge / (4 * math.Pi * Eps0 * distanceSquared)
	F := Fg - Fc

	direction := deltaPosition.Normalize()

	return direction.Mul(F)
}
