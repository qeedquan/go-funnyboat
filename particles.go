package main

import (
	"math/rand"

	"github.com/qeedquan/go-media/sdl"
	"github.com/qeedquan/go-media/sdl/sdlgfx"
	"github.com/qeedquan/go-media/sdl/sdlimage/sdlcolor"
)

type Particle struct {
	Pos        Point
	Vel        Point
	Color      sdl.Color
	Accel      Point
	Size       int
	Initial    int
	Life       int
	Opacity    float64
	Underwater bool
	image      *Image
}

func (p *Particle) Update() {
	p.Pos = p.Pos.Add(p.Vel)
	p.Vel = p.Vel.Add(p.Accel)
	if p.Life > 0 {
		p.Life--
	}

	if !p.Underwater && p.Vel.Y > 0 {
		p.Life = 0
	}
}

func (p *Particle) Draw() {
	p.image.Bind()

	screen.SetDrawColor(sdlcolor.Transparent)
	screen.Clear()

	x := p.Size / 2
	y := x
	r := p.Size / 2
	p.image.SetAlphaMod(uint8(float64(p.Life) * 255 * p.Opacity / float64(p.Initial)))
	sdlgfx.FilledEllipse(screen.Renderer, x, y, r, r, p.Color)

	p.image.Unbind()
	p.image.Blit(p.Pos)
}

func (p *Particle) Free() {
	p.image.Destroy()
}

type Ensemble struct {
	image     *Image
	particles []Particle
}

func (e *Ensemble) Blood(pos Point) {
	p := Particle{
		Pos:        pos,
		Vel:        Point{rand.Float64()*5 - 2.5, rand.Float64()*5 - 2.5},
		Color:      sdl.Color{230, 30, 20, 255},
		Accel:      Point{0, 0.7},
		Size:       rand.Intn(5) + 1,
		Initial:    rand.Intn(30),
		Opacity:    1,
		Underwater: true,
	}
	p.image = NewImage(p.Size, p.Size)
	p.Life = p.Initial
	e.particles = append(e.particles, p)
}

func (e *Ensemble) Explosion(pos Point) {
	p := Particle{
		Pos:        pos,
		Vel:        Point{rand.Float64()*5 - 2.5, rand.Float64()*5 - 2.5},
		Color:      sdl.Color{230, 30 + uint8(rand.Intn(200)), 20, 255},
		Accel:      Point{0, 0.2},
		Size:       rand.Intn(7) + 1,
		Initial:    rand.Intn(30),
		Opacity:    1,
		Underwater: true,
	}
	p.image = NewImage(p.Size, p.Size)
	p.Life = p.Initial
	e.particles = append(e.particles, p)
}

func (e *Ensemble) Water(pos Point) {
	p := Particle{
		Pos:        pos,
		Vel:        Point{rand.Float64()*5 - 2.5, -rand.Float64()*2.5 - 2},
		Color:      sdl.Color{20, 60, 180, 255},
		Accel:      Point{0, 0.3},
		Size:       rand.Intn(5) + 1,
		Initial:    rand.Intn(30),
		Opacity:    0.5,
		Underwater: false,
	}
	p.image = NewImage(p.Size, p.Size)
	p.Life = p.Initial
	e.particles = append(e.particles, p)
}

func (e *Ensemble) Debris(pos Point) {
	p := Particle{
		Pos:        pos,
		Vel:        Point{rand.Float64()*5 - 2.5, rand.Float64()*5 - 2.5},
		Color:      sdl.Color{90, 90, 90, 255},
		Accel:      Point{0, 0.2},
		Size:       rand.Intn(7) + 1,
		Initial:    rand.Intn(30),
		Opacity:    1,
		Underwater: true,
	}
	p.image = NewImage(p.Size, p.Size)
	p.Life = p.Initial
	e.particles = append(e.particles, p)
}

func (e *Ensemble) Wood(pos Point) {
	p := Particle{
		Pos:        pos,
		Vel:        Point{rand.Float64()*5 - 2.5, rand.Float64()*5 - 2.5},
		Color:      sdl.Color{148, 69, 6, 255},
		Accel:      Point{0, 0.2},
		Size:       rand.Intn(7) + 1,
		Initial:    rand.Intn(30),
		Opacity:    1,
		Underwater: true,
	}
	p.image = NewImage(p.Size, p.Size)
	p.Life = p.Initial
	e.particles = append(e.particles, p)
}

func (e *Ensemble) Steam(pos Point) {
	p := Particle{
		Pos:        pos,
		Vel:        Point{-rand.Float64() * 0.3, -rand.Float64() * 0.1},
		Color:      sdl.Color{240, 240, 240, 255},
		Accel:      Point{-0.1, -0.00002},
		Size:       rand.Intn(10) + 1,
		Initial:    rand.Intn(30),
		Opacity:    0.5,
		Underwater: true,
	}
	p.image = NewImage(p.Size, p.Size)
	p.Life = p.Initial
	e.particles = append(e.particles, p)
}

func (e *Ensemble) Fire(pos Point) {
	p := Particle{
		Pos:        pos,
		Vel:        Point{-rand.Float64() * 0.3, -rand.Float64() * 0.1},
		Color:      sdl.Color{255, 210, 170, 255},
		Accel:      Point{-0.1, -0.00002},
		Size:       rand.Intn(11) + 1,
		Initial:    rand.Intn(30),
		Opacity:    0.4,
		Underwater: false,
	}
	p.image = NewImage(p.Size, p.Size)
	p.Life = p.Initial
	e.particles = append(e.particles, p)
}

func (e *Ensemble) Trace(pos Point) {
	p := Particle{
		Pos:        pos,
		Vel:        Point{},
		Color:      sdl.Color{170, 170, 170, 255},
		Accel:      Point{},
		Size:       6,
		Initial:    5 + rand.Intn(5),
		Opacity:    0.1 + rand.Float64()*0.1,
		Underwater: false,
	}
	p.image = NewImage(p.Size, p.Size)
	p.Life = p.Initial
	e.particles = append(e.particles, p)
}

func (e *Ensemble) Init() {
	e.image = NewImage(W, H)
}

func (e *Ensemble) Update() {
	for i := 0; i < len(e.particles); {
		p := &e.particles[i]
		p.Update()
		if p.Life <= 0 {
			p.Free()
			l := len(e.particles) - 1
			e.particles[i], e.particles = e.particles[l], e.particles[:l]
		} else {
			i++
		}
	}
}

func (e *Ensemble) Draw() {
	if !config.Particles {
		return
	}

	for i := range e.particles {
		p := &e.particles[i]
		p.Draw()
	}
}

func (e *Ensemble) Free() {
	for _, p := range e.particles {
		p.Free()
	}
	e.particles = e.particles[:0]
}
