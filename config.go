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

func constructSystem(filepath string) ([]Particle, []float64, []uint32) {
	var c Config
	if _, err := toml.DecodeFile(filepath, &c); err != nil {
		log.Fatal(err)
	}
	var rp []Particle
	var rr []float64
	var textures []uint32
	for _, b := range c.Bodies {
		t := Particle{mgl64.Vec3{b.Distance, 0, 0}, mgl64.Vec3{0, 0, b.Speed}, b.Mass, 0}
		rp = append(rp, t)
		rr = append(rr, b.Diameter/2)
		text, err := newTexture("textures/" + b.Texture)
		fmt.Println(text)
		if err != nil {
			log.Fatal(err)
		}
		textures = append(textures, text)
	}
	return rp, rr, textures
}
