package main

import (
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

type Surface [3]uint32

type Mesh struct {
	Vertices []mgl64.Vec3
	Normals  []mgl64.Vec3
	UVcoords []mgl64.Vec2
	Faces    []Surface
}

func (m *Mesh) Enhance() {
	adj := make([][]int, len(m.Vertices))
	for i := range adj {
		adj[i] = make([]int, i+1)
		for j := range adj[i] {
			adj[i][j] = -1
		}
	}

	for i, f := range m.Faces {
		for j := range f {
			a := f[j]
			b := f[(j+1)%3]
			if a < b {
				a, b = b, a
			}

			if adj[a][b] == -1 {
				adj[a][b] = len(m.Vertices)
				m.Vertices = append(m.Vertices, lerp64(m.Vertices[a], m.Vertices[b], 0.5))
				//m.Colors = append(m.Colors, lerp32(m.Colors[a], m.Colors[b], 0.5))
				//m.Colors = append(m.Colors, mgl32.Vec3{rand.Float32(), rand.Float32(), rand.Float32()})
			}
			m.Faces[i][j] = uint32(adj[a][b])
		}
		m.Faces = append(m.Faces,
			Surface{f[0], m.Faces[i][0], m.Faces[i][2]},
			Surface{f[1], m.Faces[i][1], m.Faces[i][0]},
			Surface{f[2], m.Faces[i][2], m.Faces[i][1]},
		)
	}
}

func (m *Mesh) PuffUp(radius float64) {
	m.UVcoords = make([]mgl64.Vec2, len(m.Vertices))
	for i, p := range m.Vertices {
		_, theta, phi := mgl64.CartesianToSpherical(p)
		m.Vertices[i] = mgl64.SphericalToCartesian(radius, theta, phi)
		m.UVcoords[i] = mgl64.Vec2{(phi + math.Pi), 2 * theta}.Mul(1 / (2 * math.Pi))
	}
}
