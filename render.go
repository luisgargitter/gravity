package main

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	"log"
	"os"
	"unsafe"

	"github.com/go-gl/gl/v2.1/gl"
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
	texture   uint32
	Vao       VAO
}

func ConstructVBO(vertices []mgl64.Vec3, uvcoords []mgl64.Vec2) VBO {
	if len(vertices) != len(uvcoords) {
		log.Fatal("mismatch in amount of vertives and uvcoords")
	}

	var r uint32
	a := make([][5]float64, len(uvcoords))
	for i, uv := range uvcoords {
		v := vertices[i]
		a[i] = [5]float64{v[0], v[1], v[2], uv[0], uv[1]}
	}

	gl.GenBuffers(1, &r)
	gl.BindBuffer(gl.ARRAY_BUFFER, r)
	gl.BufferData(
		gl.ARRAY_BUFFER,
		int(unsafe.Sizeof(a[0]))*len(a),
		unsafe.Pointer(&a[0]),
		gl.STATIC_DRAW,
	)

	return VBO(r)
}

func ConstructEBO(faces []Surface) EBO {
	var r uint32
	gl.GenBuffers(1, &r)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, r)
	gl.BufferData(
		gl.ELEMENT_ARRAY_BUFFER,
		int(unsafe.Sizeof(faces[0]))*len(faces),
		unsafe.Pointer(&faces[0]),
		gl.STATIC_DRAW,
	)

	return EBO{r, int32(len(faces))}
}

func ConstructVAO(vbo VBO, ebo EBO) VAO {
	var r uint32
	gl.GenVertexArrays(1, &r)
	gl.BindVertexArray(r)

	gl.BindBuffer(gl.ARRAY_BUFFER, uint32(vbo))
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo.ebo)

	gl.VertexAttribPointer(0, 3, gl.DOUBLE, false, int32(unsafe.Sizeof([5]float64{})), nil)
	gl.VertexAttribPointer(1, 2, gl.DOUBLE, false, int32(unsafe.Sizeof([5]float64{})), unsafe.Pointer(unsafe.Sizeof([3]float64{})))
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)

	gl.BindVertexArray(0)

	return VAO{r, ebo.count}
}

func (m *Mesh) Load() VAO {
	vbo := ConstructVBO(m.Vertices, m.UVcoords)
	ebo := ConstructEBO(m.Faces)

	return ConstructVAO(vbo, ebo)
}

func newTexture(file string) (uint32, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return 0, fmt.Errorf("texture %q not found on disk: %v", file, err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return texture, nil
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

		gl.BindTexture(gl.TEXTURE0, o.texture)
		o.Vao.Draw()
	}
}
