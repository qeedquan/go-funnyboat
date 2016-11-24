package main

import "math"

type Shark struct {
	Entity
	Step int
}

func (s *Shark) Init() {
	p := LoadImage("hai")
	i := p.Copy()

	s.Entity.Reset()
	s.Pictures = []*Image{p}
	s.Images = []*Image{i}
	s.Sound = LoadSound("kraah")
	s.Pos = Point{W, 0}
	s.Life = 1
}

func (s *Shark) Update() {
	m := s.Image()
	l := WaterRegion(s.Pos, m)
	if s.Dying {
		s.UpdateAngle(Lerp(s.Angle, s.TargetAngle, 0.6))
		s.Pos = s.Pos.Add(s.Vel)
		s.Vel.Y++

		if s.Pos.Bottom(m) > l[1] {
			s.Vel = s.Vel.Scale(0.8)
		}

		if s.Pos.Y >= H {
			s.Dead = true
		}

		return
	}

	if !s.Jumping {
		s.Pos.Y = l[1] - 8
		s.TargetAngle = math.Atan((l[0]-l[2])/32) * Degree
	} else {
		s.Vel.Y++
		if s.Pos.Y > l[1]-8 {
			s.Jumping = false
			s.Vel = Point{-2, 0}
		}
	}

	if s.Step%40 == 0 {
		s.Jumping = true
		s.Vel = Point{-3, -10}
	}

	s.Pos = s.Pos.Add(s.Vel)
	s.Step++
	s.UpdateAngle(Lerp(s.Angle, s.TargetAngle, 0.8))
}
