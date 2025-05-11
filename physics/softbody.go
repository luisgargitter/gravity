package physics

import (
	"github.com/go-gl/mathgl/mgl64"
)

type Softbody = SimpleUndirectedGraph[Particle, Link]

func softbodyAdd(d *Softbody, a *Softbody, b *Softbody) *Softbody {
	for i := range a.vertices {
		particleAdd(&d.vertices[i], &a.vertices[i], &b.vertices[i])
	}
	for i := range a.edges {
		linkAdd(&d.edges[i].weight, &a.edges[i].weight, &b.edges[i].weight)
	}
	return d
}

func softbodyMul(d *Softbody, a *Softbody, c float64) *Softbody {
	for i := range a.vertices {
		particleMul(&d.vertices[i], &a.vertices[i], c)
	}
	for i := range a.edges {
		linkMul(&d.edges[i].weight, &a.edges[i].weight, c)
	}
	return d
}

func dSoftbody(y *Softbody, dy *Softbody) {
	for i := range dy.vertices {
		dy.vertices[i].Velocity = mgl64.Vec3{0, 0, 0}
	}

	for i := range y.edges {
		e := &y.edges[i]
		p1 := &y.vertices[e.start]
		p2 := &y.vertices[e.end]

		f := p1.DampenedSpringForceV(p2, &e.weight)
		fp1 := f.Mul(1.0 / p1.Mass)
		fp2 := f.Mul(-1.0 / p2.Mass)

		dy.vertices[e.start].Velocity = dy.vertices[e.start].Velocity.Add(fp1)
		dy.vertices[e.end].Velocity = dy.vertices[e.end].Velocity.Add(fp2)

		dy.edges[i].weight = Link{0, 0, 0}
	}
	for i := range dy.vertices {
		dy.vertices[i].Position = y.vertices[i].Velocity
		dy.vertices[i].Mass = 0
		dy.vertices[i].Charge = 0
	}
}
