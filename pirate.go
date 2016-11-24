package main

type Pirate struct {
	Boat
}

func (p *Pirate) Init() {
	m := LoadImage("merkkari")
	i := m.Copy()

	p.Entity.Reset()
	p.Pictures = []*Image{m}
	p.Images = []*Image{i}
	p.Sound = LoadSound("blub")
	p.Life = 2
	p.Pos = Point{W, WaterLevel(W)}
	p.Vel = Point{-1, 0}
}

func (p *Pirate) Update() {
	UpdateEnemyBoat(&p.Entity, 0.9, 2, 0.25, 1, false)
}
