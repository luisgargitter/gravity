package main

import (
	"github.com/go-gl/mathgl/mgl64"
)

const (
	offsetPos    = 0
	offsetVel    = 3
	offsetMass   = 6
	offsetCharge = 7
	stride       = offsetCharge + 1
)

func VecNGetVec3(vn *mgl64.VecN, i int) mgl64.Vec3 {
	raw := vn.Raw()
	return mgl64.Vec3{raw[i+0], raw[i+1], raw[i+2]}
}

func VecNSetVec3(vn *mgl64.VecN, i int, v mgl64.Vec3) {
	raw := vn.Raw()
	raw[i+0], raw[i+1], raw[i+2] = v[0], v[1], v[2]
}

func VecNSetParticle(vn *mgl64.VecN, i int, p *Particle) {
	VecNSetVec3(vn, i+offsetPos, p.Position)
	VecNSetVec3(vn, i+offsetVel, p.Velocity)
	vn.Set(i+offsetMass, p.Mass)
	vn.Set(i+offsetCharge, p.Charge)
}

func VecNGetParticle(vn *mgl64.VecN, i int) *Particle {
	return &Particle{
		VecNGetVec3(vn, i+offsetPos),
		VecNGetVec3(vn, i+offsetVel),
		vn.Get(i + offsetMass),
		vn.Get(i + offsetCharge),
	}
}

func ParticlesToVecN(particles []Particle) *mgl64.VecN {
	vn := mgl64.NewVecN(len(particles) * stride)
	for i := range particles {
		VecNSetParticle(vn, i*stride, &particles[i])
	}
	return vn
}

func VecNToParticles(vn *mgl64.VecN) []Particle {
	particles := make([]Particle, vn.Size()/stride)
	for i := range particles {
		particles[i] = *VecNGetParticle(vn, i*stride)
	}
	return particles
}

func dParticleSystem(y *mgl64.VecN, dy *mgl64.VecN) {
	for i := 0; i < y.Size(); i += stride {
		VecNSetVec3(dy, i+offsetVel, mgl64.Vec3{0, 0, 0})
	}
	for i := 0; i < y.Size(); i += stride {
		p1 := VecNGetParticle(y, i)
		for j := i + stride; j < y.Size(); j += stride {
			p2 := VecNGetParticle(y, j)

			f := p1.ForceV(p2)
			fp1 := f.Mul(1.0 / p1.Mass)
			fp2 := f.Mul(-1.0 / p2.Mass)

			VecNSetVec3(dy, i+offsetVel, fp1.Add(VecNGetVec3(dy, i+offsetVel)))
			VecNSetVec3(dy, j+offsetVel, fp2.Add(VecNGetVec3(dy, j+offsetVel)))
		}
		// change in Position
		VecNSetVec3(dy, i+offsetPos, VecNGetVec3(y, i+offsetVel))
		// ensure mass and charge do not change
		dy.Set(i+offsetMass, 0)
		dy.Set(i+offsetCharge, 0)
	}
}
