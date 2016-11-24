package main

import (
	"math/rand"

	"github.com/qeedquan/go-media/sdl"
)

type Mine struct {
	Entity
	Exploding     bool
	ExplodeFrames int
}

func (m *Mine) Init() {
	h := H - int(WaterLevel(rand.Float64()*320)) - 4
	p := LoadImage("miina")
	i := p.CopySize(p.W, p.H+h)

	m.Entity.Reset()
	m.Pictures = []*Image{p}
	m.Images = []*Image{i}
	m.Frame = 0
	m.Sound = LoadSound("poks")
	m.Pos = Point{W, H - float64(i.H)}
	m.Vel = Point{-1, 0}

	i.Bind()
	screen.SetDrawColor(sdl.Color{25, 25, 25, 255})
	x := i.W / 2
	i.Vline(x, i.H, p.H)
	i.Unbind()
}

func (m *Mine) Update() {
	m.Pos = m.Pos.Add(m.Vel)

	i := m.Image()
	l := WaterLevel(m.Pos.CenterX(i))
	if H-l < float64(i.H)-4 {
		m.Pos.Y = l - 4
	} else {
		m.Pos.Y = H - float64(i.H)
	}

	if m.Exploding {
		if m.ExplodeFrames > 0 {
			m.ExplodeFrames--
		}
	}
}

func (m *Mine) Explode() {
	m.Sound.Play(0)
	m.Exploding = true
	m.ExplodeFrames = 10
}
