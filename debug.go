package main

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl64"
)

type Info struct {
	Position    *mgl64.Vec3
	Orientation *mgl64.Vec3
	CpuStart    *float64
	CpuEnd      *float64
	GpuStart    *float64
	GpuEnd      *float64
}

func (i *Info) Print() {
	cpuTime := (*i.CpuEnd - *i.CpuStart)
	gpuTime := (*i.GpuEnd - *i.GpuStart)
	fps := 1 / (cpuTime + gpuTime)

	fmt.Print("\033[H\033[2J") //clears the screen
	fmt.Printf(
		"Position: (%.2f, %.2f, %.2f) Orientation: (%.2f, %.2f, %.2f) CPU: %.2f ms, GPU: %.2f ms, FPS: %.2f ",
		i.Position[0], i.Position[1], i.Position[2],
		i.Orientation[0], i.Orientation[1], i.Orientation[2],
		cpuTime*1000, gpuTime*1000, fps,
	)
}
