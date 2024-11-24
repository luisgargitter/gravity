package main

import "github.com/go-gl/mathgl/mgl64"

type Vectorize interface {
	toVecN(vn *mgl64.VecN, i int)
	fromVecN(vn *mgl64.VecN, i int)
}

type Vec3 struct {
	v *mgl64.Vec3
}

func (v *Vec3) toVecN(vn *mgl64.VecN, i int) {
	raw := vn.Raw()
	raw[i+0], raw[i+1], raw[i+2] = v.v[0], v.v[1], v.v[2]
}

func (v *Vec3) fromVecN(vn *mgl64.VecN, i int) {
	raw := vn.Raw()
	*v.v = mgl64.Vec3{raw[i+0], raw[i+1], raw[i+2]}
}

const (
	particleOffsetPos    = 0
	particleOffsetVel    = 3
	particleOffsetMass   = 6
	particleOffsetCharge = 7
	particleStride       = 8
)

func (p *Particle) toVecN(vn *mgl64.VecN, i int) {
	(&Vec3{&p.Position}).toVecN(vn, i+particleOffsetPos)
	(&Vec3{&p.Velocity}).toVecN(vn, i+particleOffsetVel)
	vn.Set(i+particleOffsetMass, p.Mass)
	vn.Set(i+particleOffsetCharge, p.Charge)
}

func (p *Particle) fromVecN(vn *mgl64.VecN, i int) {
	(&Vec3{&p.Position}).fromVecN(vn, i+particleOffsetPos)
	(&Vec3{&p.Velocity}).fromVecN(vn, i+particleOffsetVel)
	p.Mass = vn.Get(i + particleOffsetMass)
	p.Charge = vn.Get(i + particleOffsetCharge)
}

const (
	linkOffsetLength = 0
	linkOffsetSpring = 1
	linkOffsetDamper = 2
	linkStride       = 3
)

func (l *Link) toVecN(vn *mgl64.VecN, i int) {
	vn.Set(i+linkOffsetLength, l.length)
	vn.Set(i+linkOffsetSpring, l.springConstant)
	vn.Set(i+linkOffsetDamper, l.damperConstant)
}

func (l *Link) fromVecN(vn *mgl64.VecN, i int) {
	l.length = vn.Get(i + linkOffsetLength)
	l.springConstant = vn.Get(i + linkOffsetSpring)
	l.damperConstant = vn.Get(i + linkOffsetDamper)
}

func (e *Edge[V]) toVecN(vn *mgl64.VecN, i int) {}
