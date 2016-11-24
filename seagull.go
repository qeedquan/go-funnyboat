package main

import (
	"fmt"
	"math/rand"
)

type Seagull struct {
	Entity
	Step uint
}

func (s *Seagull) Init() {
	var p, i []*Image

	for n := 1; n <= 3; n++ {
		p = append(p, LoadImage(fmt.Sprint("lokki", n)))
		i = append(i, p[n-1].Copy())
	}

	s.Entity.Reset()
	s.Pictures = p
	s.Images = i
	s.Pos = Point{W, H/10 + rand.Float64()*H/10}
	s.Vel = Point{-2, 0}
	s.Life = 1
}

func (s *Seagull) Update() {
	s.Step++

	if !s.Dying {
		if s.Step%3 == 0 {
			s.Frame = (s.Frame + 1) % len(s.Images)
		}
	}

	s.Pos = s.Pos.Add(s.Vel)
	s.UpdateAngle(Lerp(s.Angle, s.TargetAngle, 0.2))

	if s.Dying {
		s.Vel.Y++
		if s.Pos.X < 0 || s.Pos.X >= W || s.Pos.Y < 0 || s.Pos.Y >= H {
			s.Dead = true
		}
	}
}
