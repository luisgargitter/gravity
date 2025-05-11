package graphics

import (
	"github.com/go-gl/mathgl/mgl32"
)

func Tetraheadron() Mesh {
	var r Mesh
	r.Vertices = []mgl32.Vec3{
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
	r.Vertices = []mgl32.Vec3{
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
		{0, 1, 4},
		{0, 2, 1},
		{0, 4, 2},
		{3, 1, 2},
		{3, 2, 7},
		{3, 7, 1},
		{5, 1, 7},
		{5, 4, 1},
		{5, 7, 4},
		{6, 2, 4},
		{6, 4, 7},
		{6, 7, 2},
	}
	return r
}

func Sphere(resolution int) Mesh {
	r := Cube()
	for range resolution {
		r.Enhance()
	}
	r.PuffUp()

	return r
}
