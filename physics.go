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

const (
	particleOffsetPos    = 0
	particleOffsetVel    = 3
	particleOffsetMass   = 6
	particleOffsetCharge = 7
)

func VecNSetParticle(vn *mgl64.VecN, i int, p *Particle) {
	VecNSetVec3(vn, i+particleOffsetPos, p.Position)
	VecNSetVec3(vn, i+particleOffsetVel, p.Velocity)
	vn.Set(i+particleOffsetMass, p.Mass)
	vn.Set(i+particleOffsetCharge, p.Charge)
}

func VecNGetParticle(vn *mgl64.VecN, i int) *Particle {
	return &Particle{
		VecNGetVec3(vn, i+particleOffsetPos),
		VecNGetVec3(vn, i+particleOffsetVel),
		vn.Get(i + particleOffsetMass),
		vn.Get(i + particleOffsetCharge),
	}
}

type Link struct {
	length         float64
	springConstant float64
	damperConstant float64
	start          int
	end            int // start < end
}

func VecNSetLink(vn *mgl64.VecN, i int, l *Link) {
	VecNSetVec3(vn, i, mgl64.Vec3{l.length, l.springConstant, l.damperConstant})
}

func VecNGetLink(vn *mgl64.VecN, i int) *Link {
	v := VecNGetVec3(vn, i)
	return &Link{v[0], v[1], v[2]}
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
