package graphics

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Scene struct {
	camera  *Camera
	program uint32
	meshes  map[*Mesh]*VAO
	objects []Object
}

func SceneInit(camera Camera, program uint32, objects []Object) *Scene {
	s := Scene{}
	s.camera = &camera
	s.program = program
	s.meshes = make(map[*Mesh]*VAO)
	s.objects = objects
	for _, o := range objects {
		if s.meshes[o.Mesh] == nil {
			vao := o.Mesh.Load()
			s.meshes[o.Mesh] = &vao
		}
	}
	return &s
}

func (s *Scene) Draw() {
	viewUni := gl.GetUniformLocation(s.program, gl.Str("view\x00"))
	cp := s.camera.Perspective()

	gl.ClearColor(0, 0, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	for _, o := range s.objects {
		ou := mgl32.Translate3D(o.Position[0], o.Position[1], o.Position[2]).Mul4(o.Orientation.Mat4())
		vu := cp.Mul4(ou)
		gl.UniformMatrix4fv(viewUni, 1, false, &vu[0])

		m := s.meshes[o.Mesh]
		m.Draw()
	}
}
