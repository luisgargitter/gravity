package main

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
)

type Trail struct {
	Position *mgl64.Vec3
	Length   int
	Width    float32
	Color    mgl32.Vec3
	Curve    []mgl32.Vec3
}

func (t *Trail) Update(scale float64) {
	p := (*t.Position).Mul(scale)
	t.Curve = append(t.Curve, mgl32.Vec3{float32(p[0]), float32(p[1]), float32(p[2])})
	if len(t.Curve) > t.Length {
		t.Curve = t.Curve[1:]
	}
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
