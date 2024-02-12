package main

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
)

type PointCloud []mgl64.Vec3
type Surface [3]int

type Mesh struct {
	Points PointCloud
	Colors []mgl32.Vec3
	Faces  []Surface
}

func (m *Mesh) Draw() {
	gl.Begin(gl.TRIANGLES)

	for _, f := range m.Faces {
		for _, v := range f {
			t := m.Points[v]
			gl.Color4f(m.Colors[v][0], m.Colors[v][1], m.Colors[v][2], 1)
			gl.Vertex3f(float32(t[0]), float32(t[1]), float32(t[2]))
		}
	}

	gl.End()
}

func lerp64(a, b mgl64.Vec3, t float64) mgl64.Vec3 {
	return a.Add(b.Sub(a).Mul(t))
}

func lerp32(a, b mgl32.Vec3, t float32) mgl32.Vec3 {
	return a.Add(b.Sub(a).Mul(t))
}

func (m *Mesh) Enhance() {
	adj := make([][]int, len(m.Points))
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
				adj[a][b] = len(m.Points)
				m.Points = append(m.Points, lerp64(m.Points[a], m.Points[b], 0.5))
				m.Colors = append(m.Colors, lerp32(m.Colors[a], m.Colors[b], 0.5))
				//m.Colors = append(m.Colors, mgl32.Vec3{rand.Float32(), rand.Float32(), rand.Float32()})
			}
			m.Faces[i][j] = adj[a][b]
		}
		m.Faces = append(m.Faces,
			Surface{f[0], m.Faces[i][0], m.Faces[i][2]},
			Surface{f[1], m.Faces[i][1], m.Faces[i][0]},
			Surface{f[2], m.Faces[i][2], m.Faces[i][1]},
		)
	}
}

func (m *Mesh) PuffUp(radius float64) {
	for i, p := range m.Points {
		_, theta, phi := mgl64.CartesianToSpherical(p)
		m.Points[i] = mgl64.SphericalToCartesian(radius, theta, phi)
	}
}

func Tetraheadron() Mesh {
	var r Mesh
	r.Points = []mgl64.Vec3{
		{1, 1, 1},
		{-1, -1, 1},
		{1, -1, -1},
		{-1, 1, -1},
	}
	r.Faces = []Surface{
		{0, 1, 2},
		{0, 1, 3},
		{0, 2, 3},
		{1, 2, 3},
	}

	return r
}

func Cube() Mesh {
	var r Mesh
	r.Points = []mgl64.Vec3{
		{-1, -1, -1},
		{-1, -1, 1},
		{-1, 1, -1},
		{-1, 1, 1},
		{1, -1, -1},
		{1, -1, 1},
		{1, 1, -1},
		{1, 1, 1},
	}
	r.Faces = []Surface{
		{0, 1, 2},
		{0, 1, 4},
		{0, 2, 4},
		{3, 7, 2},
		{3, 7, 1},
		{3, 2, 1},
		{5, 1, 7},
		{5, 1, 4},
		{5, 4, 7},
		{6, 2, 4},
		{6, 2, 7},
		{6, 4, 7},
	}
	return r
}
