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

func Arrayf64Tof32(s []float64, d *[]float32) {
	for i := range s {
		(*d)[i] = float32(s[i])
	}
}

func triangleNumber(n int) int {
	return n * (n - 1) / 2
}
