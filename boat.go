package main

import "math"

type Boat struct {
	Entity
}

func UpdateEnemyBoat(e *Entity, i1, i2, i3, i4 float64, stationary bool) {
	m := e.Image()
	l := WaterRegion(e.Pos, m)

	e.T++

	bottom := e.Pos.Bottom(m)

	if e.Dying {
		e.UpdateAngle(Lerp(e.Angle, e.TargetAngle, i1))
		e.Pos.Y += e.Vel.Y
		e.Vel.Y++

		if bottom > l[1] {
			e.Vel.Y *= 0.8
		}

		if e.Pos.Y >= H {
			e.Dead = true
		}

		return
	}

	if bottom > l[1]+4 {
		e.Vel.Y *= 0.8
		if e.Pos.Y > l[1] {
			e.Vel.Y -= i2
		} else {
			e.Vel.Y -= i3 * (bottom - l[1])
		}

		e.TargetAngle = i4*math.Atan((l[0]-l[2])/32)*Degree + math.Sin(float64(e.T)*0.05)*5
	}

	if stationary && e.Pos.Right(e.Image()) < W {
		e.Vel.X = 0
	}

	e.Vel.Y++
	e.Pos = e.Pos.Add(e.Vel)
	e.UpdateAngle(Lerp(e.Angle, e.TargetAngle, 0.8))
}

func (b *Boat) Die() {
	snd := b.Sound
	b.Sound = nil
	b.Entity.Die()
	b.Sound = snd
	b.Sound.Play(3)
}
