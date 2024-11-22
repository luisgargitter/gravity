package main

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
)

func lerp64(a, b mgl64.Vec3, t float64) mgl64.Vec3 {
	return a.Add(b.Sub(a).Mul(t))
}

func lerp32(a, b mgl32.Vec3, t float32) mgl32.Vec3 {
	return a.Add(b.Sub(a).Mul(t))
}

func VecNGetVec3(vn *mgl64.VecN, i int) mgl64.Vec3 {
	raw := vn.Raw()
	return mgl64.Vec3{raw[i+0], raw[i+1], raw[i+2]}
}

func VecNSetVec3(vn *mgl64.VecN, i int, v mgl64.Vec3) {
	raw := vn.Raw()
	raw[i+0], raw[i+1], raw[i+2] = v[0], v[1], v[2]
}
