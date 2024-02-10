package main

import (
	"github.com/go-gl/mathgl/mgl64"
)

const G = 0.000000000066743

type PointMass struct {
	Position mgl64.Vec3
	Inertia  mgl64.Vec3
	Mass     float64
}

func (p *PointMass) Force(a *PointMass) float64 {
	return G * p.Mass * a.Mass / a.Position.Sub(p.Position).LenSqr()
}

func (p *PointMass) ForceV(a *PointMass) mgl64.Vec3 {
	return a.Position.Sub(p.Position).Normalize().Mul(p.Force(a))
}

func (p *PointMass) Move(t float64) {
	p.Position = p.Position.Add(p.Inertia.Mul(t))
}

func (p *PointMass) ApplyForce(f mgl64.Vec3, t float64) {
	p.Inertia = p.Inertia.Add(f.Mul(t / p.Mass))
}
