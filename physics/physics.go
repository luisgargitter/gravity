package physics

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
	Length         float64
	SpringConstant float64
	DamperConstant float64
}

func linkAdd(d, a, b *Link) *Link {
	d.Length = a.Length + b.Length
	d.SpringConstant = a.SpringConstant + b.SpringConstant
	d.DamperConstant = a.DamperConstant + b.DamperConstant
	return d
}

func linkMul(d, a *Link, c float64) *Link {
	d.Length = a.Length * c
	d.SpringConstant = a.SpringConstant * c
	d.DamperConstant = a.DamperConstant * c
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
	compression := l.Length - distance
	Fs := compression * l.SpringConstant

	deltaVelocity := a.Velocity.Sub(p.Velocity)
	direction := deltaPosition.Mul(1.0 / distance)
	deltaVelocityAlongLink := deltaVelocity.Dot(direction)
	Fg := deltaVelocityAlongLink * l.DamperConstant
	F := Fs - Fg

	return direction.Mul(F)
}
