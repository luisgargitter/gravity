package graphics

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Surface [3]uint32

type Mesh struct {
	Vertices []mgl32.Vec3
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
				m.Vertices = append(m.Vertices, lerp32(m.Vertices[a], m.Vertices[b], 0.5))
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

func (m *Mesh) PuffUp() {
	for i, p := range m.Vertices {
		_, theta, phi := mgl32.CartesianToSpherical(p)
		m.Vertices[i] = mgl32.SphericalToCartesian(1.0, theta, phi)
	}
}

func (m *Mesh) Load() VAO {
	vbo := constructVBO(m.Vertices)
	ebo := constructEBO(m.Faces)

	return constructVAO(vbo, ebo)
}
