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

func particleAdd(d, a, b *Particle) *Particle {
	d.Position = a.Position.Add(b.Position)
	d.Velocity = a.Velocity.Add(b.Velocity)
	d.Mass = a.Mass + b.Mass
	d.Charge = a.Charge + b.Charge
	return d
}

func particleMul(d, a *Particle, c float64) *Particle {
	d.Position = a.Position.Mul(c)
	d.Velocity = a.Velocity.Mul(c)
	d.Mass = a.Mass * c
	d.Charge = a.Charge * c
	return d
}

type Link struct {
	length         float64
	springConstant float64
	damperConstant float64
}

func linkAdd(d, a, b *Link) *Link {
	d.length = a.length + b.length
	d.springConstant = a.springConstant + b.springConstant
	d.damperConstant = a.damperConstant + b.damperConstant
	return d
}

func linkMul(d, a *Link, c float64) *Link {
	d.length = a.length * c
	d.springConstant = a.springConstant * c
	d.damperConstant = a.damperConstant * c
	return d
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

func (p *Particle) DampenedSpringForceV(a *Particle, l *Link) mgl64.Vec3 {
	deltaPosition := a.Position.Sub(p.Position)
	distance := deltaPosition.Len()
	compression := l.length - distance
	Fs := compression * l.springConstant

	deltaVelocity := a.Velocity.Sub(p.Velocity)
	direction := deltaPosition.Mul(1.0 / distance)
	deltaVelocityAlongLink := deltaVelocity.Dot(direction)
	Fg := deltaVelocityAlongLink * l.damperConstant
	F := Fs - Fg

	return direction.Mul(F)
}
