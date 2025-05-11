package graphics

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	Position    mgl32.Vec3
	Orientation mgl32.Vec3
	Up          mgl32.Vec3
	FovY        float32
	Aspect      float32
	Near        float32
	Far         float32
}

func (c *Camera) Perspective() mgl32.Mat4 {
	return mgl32.Perspective(c.FovY, c.Aspect, c.Near, c.Far).Mul4(mgl32.LookAtV(c.Position, c.Position.Add(c.Orientation), c.Up))
}
