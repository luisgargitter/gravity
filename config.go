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
}

type Config struct {
	Bodies []Celestialbody
}

func constructSystem(filepath string) ([]PointMass, []mgl32.Vec3) {
	var c Config
	if _, err := toml.DecodeFile(filepath, &c); err != nil {
		log.Fatal(err)
	}
	var rp []PointMass
	var rc []mgl32.Vec3
	for _, b := range c.Bodies {
		t := PointMass{mgl64.Vec3{b.Distance, 0, 0}, mgl64.Vec3{0, 0, b.Speed}, b.Mass}
		rp = append(rp, t)
		rc = append(rc, mgl32.Vec3{b.Color[0], b.Color[1], b.Color[2]})
	}
	return rp, rc
}
