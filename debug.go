package main

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl64"
)

type Info struct {
	Inertia     *mgl64.Vec3
	Orientation *mgl64.Vec3
	CpuStart    *float64
	CpuEnd      *float64
	GpuStart    *float64
	GpuEnd      *float64
	Planets     *[]string
	Locked      *bool
	PlanetIndex *int
}

func (i *Info) Print() {
	cpuTime := (*i.CpuEnd - *i.CpuStart)
	gpuTime := (*i.GpuEnd - *i.GpuStart)
	fps := 1 / (cpuTime + gpuTime)

	locked := "none"
	if *i.Locked {
		locked = (*i.Planets)[*i.PlanetIndex]
	}

	fmt.Print("\033[H\033[2J") //clears the screen
	fmt.Printf(
		"Inertia: (%.2f, %.2f, %.2f), Orientation: (%.2f, %.2f, %.2f), Locked: %s, CPU: %.2f ms, GPU: %.2f ms, FPS: %.2f ",
		i.Inertia[0], i.Inertia[1], i.Inertia[2],
		i.Orientation[0], i.Orientation[1], i.Orientation[2],
		locked,
		cpuTime*1000, gpuTime*1000, fps,
	)
}
