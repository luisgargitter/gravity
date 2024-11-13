package main

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Trail struct {
	Position *mgl32.Vec3
	Length   int
	Width    float32
	Color    mgl32.Vec3
	Curve    []mgl32.Vec3
}

func NewTrail(position *mgl32.Vec3, length int, width float32, color mgl32.Vec3) *Trail {
	return &Trail{position, length, width, color, make([]mgl32.Vec3, length)}
}

func (t *Trail) Update() {
	p := *t.Position
	for i := 0; i < len(t.Curve)-1; i++ {
		t.Curve[i+1] = t.Curve[i]
	}
	t.Curve[0] = p
}

func (t *Trail) Draw() {
	gl.LineWidth(t.Width)
	gl.Begin(gl.LINE_STRIP)

	for i := 0; i < len(t.Curve); i++ {
		fade := float64(i) / float64(t.Length)
		gl.Color4f(t.Color[0], t.Color[1], t.Color[2], float32(fade))
		gl.Vertex3f(t.Curve[i][0], t.Curve[i][1], t.Curve[i][2])
	}

	gl.End()
}
