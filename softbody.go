package main

import (
	"github.com/go-gl/mathgl/mgl64"
)

type Softbody struct {
	particles ParticleSystem
	links     []Link
}

const linkStride = 3

func (s *Softbody) toVecN(d *mgl64.VecN) *mgl64.VecN {
	adjacenyTableSize := len(s.particles) * (len(s.particles) - 1) / 2
	if d == nil || d.Size() < len(s.particles)*particleStride+adjacenyTableSize*linkStride {
		d = mgl64.NewVecN(len(s.particles)*particleStride + adjacenyTableSize*linkStride)
	}
	d = d.Sub(d, d)
	d = s.particles.toVecN(d)

	for i := range s.links {
		j := len(s.particles) * particleStride
		l := &s.links[i]
		VecNSetLink(d, j+l.end*len(s.particles)*linkStride+l.start, l)
	}

	return d
}
