package main

import (
	"fmt"
	"math/rand"
)

type Cloud struct {
	image *Image
	pos   Point
	vel   Point
}

type Skies struct {
	clouds []Cloud
	t      int
}

var (
	skies Skies
)

func (c *Cloud) Init() {
	c.image = LoadImage(fmt.Sprint("cloud", rand.Intn(4)+1))
	c.pos = Point{W, rand.Float64() * 70}
	c.vel = Point{-1, 0}
}

func (c *Cloud) Update() {
	c.pos = c.pos.Add(c.vel)
}

func (c *Cloud) Draw() {
	c.image.Blit(c.pos)
}

func InitClouds() {
	skies.Init()
}

func UpdateClouds() {
	skies.Update()
}

func DrawClouds() {
	skies.Draw()
}

func (s *Skies) Init() {
	for i := 1; i <= 4; i++ {
		LoadImage(fmt.Sprint("cloud", i))
	}
}

func (s *Skies) Update() {
	if s.t%150 == 0 {
		var c Cloud

		c.Init()
		s.clouds = append(s.clouds, c)
	}

	for i := 0; i < len(s.clouds); {
		c := &s.clouds[i]
		c.Update()
		if c.pos.Right(c.image) < 0 {
			l := len(s.clouds) - 1
			s.clouds[i], s.clouds = s.clouds[l], s.clouds[:l]
		} else {
			i++
		}
	}

	s.t++
}

func (s *Skies) Draw() {
	for _, c := range s.clouds {
		c.Draw()
	}
}
