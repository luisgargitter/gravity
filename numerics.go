package main

import "github.com/go-gl/mathgl/mgl32"

type T = mgl32.Vec3

type ODE_F func(T, float32) T

func runge_kutta_3(f ODE_F, y0 T, t float32) T {
	k1 := f(y0, t)
	k2 := f(y0.Add(k1.Mul(0.5)), t+t*0.5)
	k3 := f(y0.Sub(k1.Mul(t).Add(k2.Mul(2.0*t))), 2.0*t)
	y1 := y0.Add(k1.Mul(1.0 / 6.0).Add(k2.Mul(4.0 / 6.0).Add(k3.Mul(1.0 / 6.0)))).Mul(t)
	return y1
}
