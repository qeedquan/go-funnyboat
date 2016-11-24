package main

type Powerup struct {
	Entity
	Picked bool
	Fading bool
	Fade   int
}

func (p *Powerup) Init() {
	p.Pictures = []*Image{LoadImage("sydan")}
	p.Images = []*Image{p.Pictures[0].Copy()}
	p.Pos = Point{W, WaterLevel(W)}
	p.Vel = Point{-1, 0}
	p.Picked = false
	p.Fading = false
}

func (p *Powerup) Update() {
	m := p.Image()
	l := WaterLevel(p.Pos.CenterX(m))
	if p.Fading {
		if p.Fade > 0 {
			p.Fade--
		} else {
			p.Picked = true
		}
	}

	bottom := p.Pos.Bottom(m)
	if bottom > l {
		p.Vel.Y *= 0.8
		if p.Pos.X > l {
			p.Vel.Y -= 2
		} else {
			p.Vel.Y -= 0.25 * (bottom - l)
		}
	}

	p.Vel.Y++
	p.Pos = p.Pos.Add(p.Vel)
}

func (p *Powerup) Pickup() {
	p.Fading = true
	p.Fade = 15
}
