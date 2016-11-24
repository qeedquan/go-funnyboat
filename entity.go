package main

import "math"

type Entity struct {
	Sound       *Sound
	Pictures    []*Image
	Images      []*Image
	Frame       int
	Pos         Point
	Vel         Point
	Accel       Point
	Life        int
	Dying       bool
	Dead        bool
	Jumping     bool
	Angle       float64
	TargetAngle float64
	T           int
}

func (e *Entity) Reset() {
	e.Frame = 0
	e.Pos = Point{}
	e.Vel = Point{}
	e.Accel = Point{}
	e.Life = 0
	e.Dying = false
	e.Dead = false
	e.Jumping = false
	e.Angle = 0
	e.TargetAngle = 0
	e.T = 0
}

func (e *Entity) Image() *Image {
	return e.Images[e.Frame]
}

func (e *Entity) Draw() {
	m := e.Image()
	m.Blit(e.Pos)
}

func (e *Entity) UpdateAngle(angle float64) {
	e.Angle = angle
	m := e.Image()
	m.UpdateAngle(angle)
}

func (e *Entity) Free() {
	for i := range e.Images {
		e.Images[i].Destroy()
	}
	e.Sound.Free()
}

func (e *Entity) Rotate(p Point) Point {
	m := e.Image()
	x := p.X - float64(m.W)/2
	y := p.Y - float64(m.H)/2
	t := -e.Angle * Radian
	return Point{
		-y*math.Sin(t) + x*math.Cos(t),
		y*math.Cos(t) + x*math.Sin(t),
	}
}

func (e *Entity) Die() {
	e.Dying = true
	e.Sound.Play(0)
	e.Vel.Y = -5
	e.TargetAngle = 90
}

func (e *Entity) Damage(dec int) {
	if e.Life -= dec; e.Life <= 0 {
		e.Die()
	}
}
