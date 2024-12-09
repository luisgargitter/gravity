package main

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl64"
	"log"
)

type Plotter struct {
	win        *glfw.Window
	width      int
	height     int
	length     int
	dataPoints []*float64
	maxs       []float64
	mins       []float64
	tracks     [][]float64
}

func InitPlotter(width, height, length int) *Plotter {
	var p Plotter
	p.width = width
	p.height = height
	p.length = length

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	win, err := glfw.CreateWindow(p.width, p.height, "Plotter", nil, nil)
	if err != nil {
		log.Fatal("error creating window for plotter: " + err.Error())
	}
	p.win = win

	p.win.MakeContextCurrent()
	glfw.SwapInterval(1)

	err = gl.Init()
	if err != nil {
		log.Fatalln("failed to initialize OpenGL", err)
	}

	p.dataPoints = make([]*float64, 0)
	p.tracks = make([][]float64, 0)
	p.maxs = make([]float64, 0)
	p.mins = make([]float64, 0)

	return &p
}

func (p *Plotter) Attach(f *float64) {
	p.dataPoints = append(p.dataPoints, f)
	p.tracks = append(p.tracks, make([]float64, p.length))
	p.maxs = append(p.maxs, 0)
	p.mins = append(p.mins, 0)
}

func (p *Plotter) Update() {
	for i := range p.dataPoints {
		j := 0
		for ; j < len(p.tracks[i])-1; j++ {
			p.tracks[i][j] = p.tracks[i][j+1]
		}
		p.tracks[i][j] = *p.dataPoints[i]
		if p.tracks[i][j] > p.maxs[i] {
			p.maxs[i] = p.tracks[i][j]
		} else if p.tracks[i][j] < p.mins[i] {
			p.mins[i] = p.tracks[i][j]
		}
	}
}

func (p *Plotter) Draw() {
	gl.ClearColor(0, 0, 0, 1)
	p.win.MakeContextCurrent()
	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.Color3f(1, 0, 0)
	for i := range p.tracks {
		gl.Begin(gl.LINE_STRIP)
		for j := range p.tracks[i] {
			gl.Vertex2d(map64(mgl64.Vec2{0, float64(p.length) - 1}, mgl64.Vec2{-1, 1}, float64(j)),
				map64(mgl64.Vec2{p.mins[i], p.maxs[i]}, mgl64.Vec2{-1, 1}, p.tracks[i][j]))
		}
		gl.End()
	}
	p.win.SwapBuffers()
}
