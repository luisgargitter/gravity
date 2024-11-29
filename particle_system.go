package main

import (
	"github.com/go-gl/mathgl/mgl64"
)

type ParticleSystem []Particle

func particleSystemAdd(d *ParticleSystem, a *ParticleSystem, b *ParticleSystem) *ParticleSystem {
	for i := range *a {
		(*d)[i].Position = (*a)[i].Position.Add((*b)[i].Position)
		(*d)[i].Velocity = (*a)[i].Velocity.Add((*b)[i].Velocity)
		(*d)[i].Mass = (*a)[i].Mass + (*b)[i].Mass
		(*d)[i].Charge = (*a)[i].Charge + (*b)[i].Charge
	}
	return d
}

func particleSystemMul(d *ParticleSystem, a *ParticleSystem, c float64) *ParticleSystem {
	for i := range *a {
		(*d)[i].Position = (*a)[i].Position.Mul(c)
		(*d)[i].Velocity = (*a)[i].Velocity.Mul(c)
		(*d)[i].Mass = (*a)[i].Mass * c
		(*d)[i].Charge = (*a)[i].Charge * c
	}
	return d
}

func dParticleSystem(y *ParticleSystem, dy *ParticleSystem) {
	for i := range *dy {
		(*dy)[i].Velocity = mgl64.Vec3{0, 0, 0}
	}

	for i := range *y {
		p1 := &(*y)[i]
		for j := i + 1; j < len(*y); j++ {
			p2 := &(*y)[j]

			f := p1.GravitationalForceV(p2)
			fp1 := f.Mul(1.0 / p1.Mass)
			fp2 := f.Mul(-1.0 / p2.Mass)

			(*dy)[i].Velocity = (*dy)[i].Velocity.Add(fp1)
			(*dy)[j].Velocity = (*dy)[j].Velocity.Add(fp2)
		}
		// change in Position
		(*dy)[i].Position = (*y)[i].Velocity
		// ensure mass and charge do not change
		(*dy)[i].Mass = 0
		(*dy)[i].Charge = 0
	}
}
