package physics

type System[T any] func(d *T, y *T)
