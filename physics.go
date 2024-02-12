package main

import (
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

const (
	G    = 6.6743015e-11
	Eps0 = 8.8541878128e-12
)

type Particle struct {
	Position mgl64.Vec3
	Inertia  mgl64.Vec3
	Mass     float64
	Charge   float64
}

func (p *Particle) Force(a *Particle) float64 {
	ds := a.Position.Sub(p.Position).LenSqr()
	Fg := G * p.Mass * a.Mass / ds
	Fc := p.Charge * a.Charge / (4 * math.Pi * Eps0 * ds)
	return Fg - Fc
}

func (p *Particle) ForceV(a *Particle) mgl64.Vec3 {
	return a.Position.Sub(p.Position).Normalize().Mul(p.Force(a))
}

func (p *Particle) Move(t float64) {
	p.Position = p.Position.Add(p.Inertia.Mul(t))
}

func (p *Particle) ApplyForce(f mgl64.Vec3, t float64) {
	p.Inertia = p.Inertia.Add(f.Mul(t / p.Mass))
}
