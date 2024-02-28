package main

import (
	"unsafe"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
)

type VBO uint32
type EBO struct {
	ebo   uint32
	count int32
}
type VAO struct {
	vao   uint32
	count int32
}

type Mesh struct {
	Points PointCloud
	Colors []mgl32.Vec3
	Faces  Faces
}

type Camera struct {
	Projection mgl64.Mat4
	POV        *Pov
}

type Scene struct {
	C  *Camera
	Os []Object
}

type Object struct {
	Transform mgl64.Mat4
	Vao       VAO
}

func (p PointCloud) Load() VBO {
	var r uint32
	gl.GenBuffers(1, &r)
	gl.BindBuffer(gl.ARRAY_BUFFER, r)
	gl.BufferData(
		gl.ARRAY_BUFFER,
		int(unsafe.Sizeof(p[0]))*len(p),
		unsafe.Pointer(&p[0]),
		gl.STATIC_DRAW,
	)

	return VBO(r)
}

func (f Faces) Load() EBO {
	var r uint32
	gl.GenBuffers(1, &r)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, r)
	gl.BufferData(
		gl.ELEMENT_ARRAY_BUFFER,
		int(unsafe.Sizeof(f[0]))*len(f),
		unsafe.Pointer(&f[0]),
		gl.STATIC_DRAW,
	)

	return EBO{r, int32(len(f))}
}

func (m *Mesh) Load() VAO {
	var r uint32
	gl.GenVertexArrays(1, &r)
	gl.BindVertexArray(r)

	vbo := m.Points.Load()
	ebo := m.Faces.Load()

	gl.BindBuffer(gl.ARRAY_BUFFER, uint32(vbo))
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo.ebo)

	gl.VertexAttribPointer(0, 3, gl.DOUBLE, false, int32(unsafe.Sizeof(mgl64.Vec3{})), nil)
	gl.EnableVertexAttribArray(0)

	gl.BindVertexArray(0)

	return VAO{r, ebo.count}
}

func (v *VAO) Draw() {
	gl.BindVertexArray(v.vao)
	gl.DrawElements(gl.TRIANGLES, int32(v.count*3), gl.UNSIGNED_INT, nil)
	gl.BindVertexArray(0)
}

func (c *Camera) ViewMatrix() mgl64.Mat4 {
	return c.Projection.Mul4(c.POV.Matrix())
}

func (s *Scene) Draw(viewUni int32) {
	m := s.C.ViewMatrix()
	for _, o := range s.Os {
		t := m.Mul4(o.Transform)
		gl.UniformMatrix4dv(viewUni, 1, false, &t[0])
		o.Vao.Draw()
	}
}
