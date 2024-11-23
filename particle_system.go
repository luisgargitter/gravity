package main

import (
	"github.com/go-gl/mathgl/mgl64"
)

type ParticleSystem []Particle

const particleStride = particleOffsetCharge + 1

func (p *ParticleSystem) toVecN(d *mgl64.VecN) *mgl64.VecN {
	if d == nil || d.Size() < len(*p)*particleStride {
		d = mgl64.NewVecN(len(*p) * particleStride)
	}

	for i := range *p {
		VecNSetParticle(d, i*particleStride, &((*p)[i]))
	}
	return d
}

func (p *ParticleSystem) fromVecN(vn *mgl64.VecN) *ParticleSystem {
	if len(*p) < vn.Size()/particleStride {
		*p = make([]Particle, vn.Size()/particleStride)
	}

	for i := range *p {
		(*p)[i] = *VecNGetParticle(vn, i*particleStride)
	}
	return p
}

func dParticleSystem(y *mgl64.VecN, dy *mgl64.VecN) {
	for i := 0; i < y.Size(); i += particleStride {
		VecNSetVec3(dy, i+particleOffsetVel, mgl64.Vec3{0, 0, 0})
	}
	for i := 0; i < y.Size(); i += particleStride {
		p1 := VecNGetParticle(y, i)
		for j := i + particleStride; j < y.Size(); j += particleStride {
			p2 := VecNGetParticle(y, j)

			f := p1.GravitationalForceV(p2)
			fp1 := f.Mul(1.0 / p1.Mass)
			fp2 := f.Mul(-1.0 / p2.Mass)

			VecNSetVec3(dy, i+particleOffsetVel, fp1.Add(VecNGetVec3(dy, i+particleOffsetVel)))
			VecNSetVec3(dy, j+particleOffsetVel, fp2.Add(VecNGetVec3(dy, j+particleOffsetVel)))
		}
		// change in Position
		VecNSetVec3(dy, i+particleOffsetPos, VecNGetVec3(y, i+particleOffsetVel))
		// ensure mass and charge do not change
		dy.Set(i+particleOffsetMass, 0)
		dy.Set(i+particleOffsetCharge, 0)
	}
}
