package main

import (
	"github.com/go-gl/mathgl/mgl64"
)

type Softbody = SimpleUndirectedGraph[Particle, Link]

const linkStride = linkOffsetDamper + 1

const (
	SoftbodyOffsetVerticesLength = 0
	SoftbodyOffsetVertices       = 1
)

func (s *Softbody) toVecN() *mgl64.VecN {
	size := 1 + len(s.vertices)*particleStride + len(s.edges)*linkStride
	d := mgl64.NewVecN(size)

	k := ParticleSystem(s.vertices)

	d.Set(SoftbodyOffsetVerticesLength, float64(len(s.vertices)))
	k.toVecN(d, SoftbodyOffsetVertices)

	startEdges := SoftbodyOffsetVertices + len(s.vertices)*particleStride
	for i, e := range s.edges {
		VecNSetLink(d, startEdges+i*linkStride, &e.weight)
	}

	return d
}

func SoftbodyFromVecN(vn *mgl64.VecN) *Softbody {
	p := ParticleSystemfromVecN(vn)

}
