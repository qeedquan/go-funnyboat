package main

import "math"

type Steamboat struct {
	Boat
	MovingLeft  bool
	MovingRight bool
	Splash      bool
	Blinks      int
}

func (s *Steamboat) Reset() {
	p := LoadImage("laiva")
	i := p.Copy()

	s.Pictures = []*Image{p}
	s.Images = []*Image{i}
	s.Sound = LoadSound("blub")

	s.Entity.Reset()
	s.Life = 5
	s.Pos = Point{50, 20}
}

func (s *Steamboat) Update() {
	m := s.Image()
	l := WaterRegion(s.Pos, m)

	s.Splash = false

	if s.Dying {
		s.UpdateAngle(Lerp(s.Angle, s.TargetAngle, 0.9))
		s.Pos.Y += s.Vel.Y
		s.Vel.Y++

		bottom := s.Pos.Bottom(m)
		if bottom > l[1] {
			s.Vel.Y *= 0.8
		}

		if s.Pos.Y >= H {
			s.Dead = true
		}

		return
	}

	s.Vel.X = 0
	switch {
	case s.MovingLeft && !s.MovingRight:
		s.Vel.X = -2
	case s.MovingRight:
		s.Vel.X = 2
	}

	bottom := s.Pos.Bottom(m)
	if bottom > l[1] {
		if s.Jumping {
			s.Splash = true
		}
		s.Jumping = false

		s.Vel.Y *= 0.8
		if s.Pos.Y > l[2] {
			s.Vel.Y -= 2
		} else {
			s.Vel.Y -= (bottom - l[1]) * 0.25
		}

		s.TargetAngle = math.Atan((l[0]-l[2])/32)*Degree + math.Sin(float64(s.T)*0.05)*5
	} else {
		s.Jumping = true
	}

	s.Vel.Y++
	s.Pos = s.Pos.Add(s.Vel)

	if s.Pos.X < 0 {
		s.Pos.X = 0
	}
	if s.Pos.Right(m) > W {
		s.Pos.X = W - float64(m.W)
	}

	s.T++
	s.UpdateAngle(Lerp(s.Angle, s.TargetAngle, 0.8))

	m.SetAlphaMod(255)
	if s.Blinks != 0 {
		if s.Blinks&1 != 0 {
			m.SetAlphaMod(0)
		}
		s.Blinks--
	}
}

func (s *Steamboat) MoveLeft(b bool) {
	if !s.Dying {
		s.MovingLeft = b
	}
}

func (s *Steamboat) MoveRight(b bool) {
	if !s.Dying {
		s.MovingRight = b
	}
}

func (s *Steamboat) Jump() {
	if !s.Dying && !s.Jumping {
		s.Jumping = true
		s.Vel.Y = -10
	}
}
