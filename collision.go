package main

import "image"

func collisionRect(e *Entity) image.Rectangle {
	x, y := int(e.Pos.X), int(e.Pos.Y)
	b := e.Image().Alpha.Bounds()
	return image.Rect(x, y, x+b.Dx(), y+b.Dy())
}

func Collision(p1, p2 *Entity) bool {
	r1 := collisionRect(p1)
	r2 := collisionRect(p2)
	r := r1.Intersect(r2)
	if r.Dx() == 0 || r.Dy() == 0 {
		return false
	}

	m1 := p1.Image().Alpha
	m2 := p2.Image().Alpha

	x1 := r.Min.X - r1.Min.X
	y1 := r.Min.Y - r1.Min.Y
	x2 := r.Min.X - r2.Min.X
	y2 := r.Min.Y - r2.Min.Y

	w, h := r.Dx(), r.Dy()
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c1 := m1.AlphaAt(x+x1, y+y1)
			c2 := m2.AlphaAt(x+x2, y+y2)
			if c1.A&c2.A != 0 {
				return true
			}
		}
	}

	return false
}
