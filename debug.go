package main

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl64"
)

type Info struct {
	Position    *mgl64.Vec3
	Inertia     *mgl64.Vec3
	Orientation *mgl64.Vec3
	CpuTime     *float64
	GpuTime     *float64
	DeltaTime   *float64
	Planets     *[]string
	Locked      *bool
	PlanetIndex *int
}

func (i *Info) Print() {
	locked := "none"
	if *i.Locked {
		locked = (*i.Planets)[*i.PlanetIndex]
	}

	fmt.Print("\033[H\033[2J") //clears the screen
	fmt.Printf(
		"Position: (%e, %e, %e), Inertia: (%e, %e, %e),	Orientation: (%f, %f, %f), Locked: %s, CPU: %.2f ms, GPU: %.2f ms, FPS: %.2f ",
		i.Position[0], i.Position[1], i.Position[2],
		i.Inertia[0], i.Inertia[1], i.Inertia[2],
		i.Orientation[0], i.Orientation[1], i.Orientation[2],
		locked,
		*i.CpuTime*1000, *i.GpuTime*1000, 1.0 / *i.DeltaTime,
	)
}
