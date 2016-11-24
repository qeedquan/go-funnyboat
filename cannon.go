package main

import "math"

type Cannon struct {
	Entity
	Special    bool
	Underwater bool
}

func (c *Cannon) Init(pos Point, shipAngle float64, left, special bool) {
	p := []*Image{LoadImage("kuti"), LoadImage("lokki2")}
	i := []*Image{p[0].Copy(), p[1].Copy()}

	c.Entity.Reset()
	c.Pictures = p
	c.Images = i
	c.Sound = LoadSound("pam")
	c.Special = special
	c.Underwater = false

	if special {
		c.Frame = 1
	}

	c.Sound.Play(0)

	angle := 0.0
	vel := 11.0
	if !left {
		if special {
			angle = -shipAngle - 15
			vel = 14
		} else {
			angle = -shipAngle - 25
		}
	} else {
		angle = -shipAngle + 180 + 25
	}

	c.Pos = pos
	c.Vel = Point{math.Cos(angle*Radian) * vel, math.Sin(angle*Radian) * vel}
}

func (c *Cannon) Update() {
	c.Pos = c.Pos.Add(c.Vel)
	c.Vel.Y += 0.4

	m := c.Image()
	if c.Pos.Bottom(m) > WaterLevel(c.Pos.CenterX(m)) {
		c.Vel.X *= 0.9
		c.Vel.Y *= 0.9
		c.Underwater = true
	}

	if c.Special && c.Vel.X != 0 {
		c.UpdateAngle(-math.Atan(c.Vel.Y/c.Vel.X) * Degree)
	}
}

func (c *Cannon) Tail() Point {
	m := c.Image()
	x := c.Pos.X
	y := c.Pos.CenterY(m) - 3 + float64(m.W)*math.Sin(c.Angle*Radian)
	return Point{x, y}
}
