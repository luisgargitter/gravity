package main

import (
	"log"

	"github.com/BurntSushi/toml"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
)

type Celestialbody struct {
	Name     string
	Color    [3]float32
	Distance float64
	Speed    float64
	Mass     float64
	Diameter float64
}

type Config struct {
	Bodies []Celestialbody
}

func constructSystem(filepath string) ([]Particle, []mgl32.Vec3, []float64) {
	var c Config
	if _, err := toml.DecodeFile(filepath, &c); err != nil {
		log.Fatal(err)
	}
	var rp []Particle
	var rc []mgl32.Vec3
	var rr []float64
	for _, b := range c.Bodies {
		t := Particle{mgl64.Vec3{b.Distance, 0, 0}, mgl64.Vec3{0, 0, b.Speed}, b.Mass, 0}
		rp = append(rp, t)
		rc = append(rc, mgl32.Vec3{b.Color[0], b.Color[1], b.Color[2]})
		rr = append(rr, b.Diameter/2)
	}
	return rp, rc, rr
}
