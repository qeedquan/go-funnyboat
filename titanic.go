package main

type Titanic struct {
	Boat
}

func (t *Titanic) Init() {
	p := LoadImage("titanic")
	i := p.Copy()

	t.Entity.Reset()
	t.Pictures = []*Image{p}
	t.Images = []*Image{i}
	t.Sound = LoadSound("blub")
	t.Life = 100
	t.Pos = Point{W, WaterLevel(W)}
	t.Pos.Y = t.Pos.Bottom(i) - t.Pos.Y
	t.Vel = Point{-1, 0}
}

func (t *Titanic) Update() {
	UpdateEnemyBoat(&t.Entity, 0.007, 1, 0.15, 0.01, true)
}
