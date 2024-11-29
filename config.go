package main

import (
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/go-gl/mathgl/mgl64"
)

type Celestialbody struct {
	Name     string
	Texture  string
	Distance float64
	Speed    float64
	Mass     float64
	Diameter float64
}

type Config struct {
	Bodies []Celestialbody
}

func constructSystem(filepath string) (ParticleSystem, []float64, []uint32, []string) {
	var c Config
	if _, err := toml.DecodeFile(filepath, &c); err != nil {
		log.Fatal(err)
	}
	var rp []Particle
	var rr []float64
	var textures []uint32
	var names []string
	for i, b := range c.Bodies {
		t := Particle{mgl64.Vec3{b.Distance, 0, 0}, mgl64.Vec3{0, 0, b.Speed}, b.Mass, 0}
		rp = append(rp, t)
		rr = append(rr, b.Diameter/2)
		text, err := newTexture("textures/" + b.Texture)
		fmt.Printf("Loading %s (%d/%d)     \r", b.Name, i, len(c.Bodies))
		if err != nil {
			log.Fatal(err)
		}
		textures = append(textures, text)
		names = append(names, b.Name)
	}
	return rp, rr, textures, names
}
