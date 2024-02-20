package main

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
)

type PointCloud []mgl64.Vec3
type Surface [3]uint32

type Mesh struct {
	Points PointCloud
	Colors []mgl32.Vec3
	Faces  []Surface
	VAO    uint32
}

type Camera struct {
	Projection  mgl64.Mat4
	Position    mgl64.Vec3
	Orientation mgl64.Vec3
}

type Scene struct {
	C  *Camera
	Os []Object
}

type Object struct {
	Transform mgl64.Mat4
	M         *Mesh
	VAO       uint32
}

func (s *Scene) Draw(viewUni int32) {
	m := mgl64.Ident4() //s.C.ViewMatrix()
	for _, o := range s.Os {
		t := m.Mul4(o.Transform)
		gl.UniformMatrix4dv(viewUni, 1, false, &t[0])
		o.Draw()
	}
}

func (o *Object) Draw() {
	gl.BindVertexArray(o.VAO)
	gl.DrawElements(gl.TRIANGLES, int32(len(o.M.Faces)*3), gl.UNSIGNED_INT, nil)
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
