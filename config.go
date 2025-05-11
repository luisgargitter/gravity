package main

import (
	"gravity/physics"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/go-gl/mathgl/mgl64"
)

type Celestialbody struct {
	Name     string
	Distance float64
	Speed    float64
	Mass     float64
	Diameter float64
}

type Config struct {
	Bodies []Celestialbody
}

func constructSystem(filepath string) (physics.ParticleSystem, []string) {
	var c Config
	if _, err := toml.DecodeFile(filepath, &c); err != nil {
		log.Fatal(err)
	}
	var rp physics.ParticleSystem
	var names []string
	for _, b := range c.Bodies {
		t := physics.Particle{Position: mgl64.Vec3{b.Distance, 0, 0}, Velocity: mgl64.Vec3{0, 0, b.Speed}, Mass: b.Mass}
		rp = append(rp, t)
		names = append(names, b.Name)
	}
	return rp, names
}
