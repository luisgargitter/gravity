package physics

import (
	"github.com/go-gl/mathgl/mgl64"
)

type ParticleSystem []Particle

func particleSystemAdd(d *ParticleSystem, a *ParticleSystem, b *ParticleSystem) *ParticleSystem {
	for i := range *a {
		particleAdd(&(*d)[i], &(*a)[i], &(*b)[i])
	}
	return d
}

func particleSystemMul(d *ParticleSystem, a *ParticleSystem, c float64) *ParticleSystem {
	for i := range *a {
		particleMul(&(*d)[i], &(*a)[i], c)
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

func ParticleSystemRK4W(particles ParticleSystem) *RK4Workspace[ParticleSystem] {
	particlesRK4W := RK4Workspace[ParticleSystem]{
		Add: particleSystemAdd,
		Mul: particleSystemMul,
		Df:  dParticleSystem,
		Y:   particles,
		D:   make(ParticleSystem, len(particles)),
		K1:  make(ParticleSystem, len(particles)),
		K2:  make(ParticleSystem, len(particles)),
		K3:  make(ParticleSystem, len(particles)),
		K4:  make(ParticleSystem, len(particles)),
	}

	return &particlesRK4W
}
