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

func (p *Particle) ForceV(a *Particle) mgl64.Vec3 {
	deltaPosition := a.Position.Sub(p.Position)
	distanceSquared := deltaPosition.LenSqr()
	Fg := G * p.Mass * a.Mass / distanceSquared
	Fc := p.Charge * a.Charge / (4 * math.Pi * Eps0 * distanceSquared)
	F := Fg - Fc

	direction := deltaPosition.Normalize()

	return direction.Mul(F)
}

func (p *Particle) Move(t float64) {
	p.Position = p.Position.Add(p.Inertia.Mul(t/p.Mass))
}

func (p *Particle) ApplyForce(f mgl64.Vec3, t float64) {
	p.Inertia = p.Inertia.Add(f.Mul(t))
}

type Link struct {
	Start int
	End int
	Length float64
	SpringCoeff float64
	DampingCoeff float64
}

type ElasticBody struct {
	Points []Particle
	Links []Link
}

func (b *ElasticBody) Step(dt float64) {
	forces := make([]mgl64.Vec3, len(b.Points))

	for i := 0; i < len(b.Links); i += 1 {
		link := b.Links[i]
		start := b.Points[link.Start]
		end := b.Points[link.End]

		deltaPosition := end.Position.Sub(start.Position)
		length := deltaPosition.Len()
		direction := deltaPosition.Mul(1/length)

		deltaVelocity := end.Inertia.Mul(1/end.Mass).Sub(start.Inertia.Mul(1/start.Mass))

		force := (link.Length - length) * link.SpringCoeff - deltaVelocity.Dot(direction) * link.DampingCoeff
		forceV := direction.Mul(force)
		forces[link.Start] = forces[link.Start].Add(forceV)
		forces[link.End] = forces[link.End].Sub(forceV)
	}

	for i := 0; i < len(b.Points); i += 1 {
		b.Points[i].ApplyForce(forces[i], dt)
	}
}