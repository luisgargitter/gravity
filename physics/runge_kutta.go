package physics

type VecAdd[T any] func(d *T, a *T, b *T) *T
type VecMul[T any] func(d *T, a *T, c float64) *T

// Runge Kutta 4
type RK4Workspace[T any] struct {
	Df  System[T]
	Add VecAdd[T]
	Mul VecMul[T]
	Y   T
	D   T
	K1  T
	K2  T
	K3  T
	K4  T
}

func RK4[T any](w *RK4Workspace[T], dt float64) {
	w.Df(&w.Y, &w.K1)
	w.Df(w.Add(&w.D, &w.Y, w.Mul(&w.D, &w.K1, dt/2.0)), &w.K2)
	w.Df(w.Add(&w.D, &w.Y, w.Mul(&w.D, &w.K2, dt/2.0)), &w.K3)
	w.Df(w.Add(&w.D, &w.Y, w.Mul(&w.D, &w.K3, dt)), &w.K4)

	// y = y0 + (h/6)(K1 + 2*K2 + 2*K3 + K4)
	w.Add(&w.D, &w.K2, &w.K3)
	w.Mul(&w.D, &w.D, 2.0)
	w.Add(&w.D, &w.D, &w.K1)
	w.Add(&w.D, &w.D, &w.K4)
	w.Mul(&w.D, &w.D, dt/6.0)
	w.Add(&w.Y, &w.Y, &w.D)
}
