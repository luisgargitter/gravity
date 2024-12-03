package main

import "github.com/go-gl/mathgl/mgl64"

type Camera struct {
	position    *mgl64.Vec3
	orientation *mgl64.Vec3
	up          *mgl64.Vec3
	fovY        float64
	aspect      float64
	near        float64
	far         float64
}

func CameraNew(position *mgl64.Vec3, orientation *mgl64.Vec3, up *mgl64.Vec3, fovY float64, aspect float64, near float64, far float64) *Camera {
	return &Camera{position, orientation, up, fovY, aspect, near, far}
}

func (c *Camera) Perspective() mgl64.Mat4 {
	return mgl64.Perspective(c.fovY, c.aspect, c.near*glCorrectionScale, c.far*glCorrectionScale).Mul4(mgl64.LookAtV(c.position.Mul(glCorrectionScale), c.position.Mul(glCorrectionScale).Add(*c.orientation), *c.up))
}
