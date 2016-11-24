package main

import (
	"math/rand"

	"github.com/qeedquan/go-media/sdl"
	"github.com/qeedquan/go-media/sdl/sdlimage/sdlcolor"
)

type Game struct {
	State

	health Health
	score  Score
	level  Level

	ensemble      Ensemble
	player        Steamboat
	playerCannons []Cannon
	enemyCannons  []Cannon
	mines         []Mine
	seagulls      []Seagull
	sharks        []Shark
	powerups      []Powerup
	pirates       []Pirate
	titanic       *Titanic

	lastShot     int
	spacePressed int
	t            int

	paused   bool
	gameOver string
}

func (g *Game) Init() {
	g.State.Init()
	g.health.Init()
	g.ensemble.Init()
}

func (g *Game) Reset(endless bool) {
	g.clear()

	g.State.Reset()
	g.health.Reset()
	g.score.Reset()
	g.level.Reset(endless)
	g.player.Reset()

	g.paused = false
	g.gameOver = ""

	g.lastShot = 0
	g.spacePressed = 0
	g.t = 0
}

func (g *Game) Run(endless bool) int {
	g.Reset(endless)

	for !g.Done {
		g.draw()
		g.update()
		g.event()
	}

	return g.score.Value
}

func (g *Game) draw() {
	g.State.Draw()
	g.health.Draw()
	g.score.Draw()

	for i := range g.powerups {
		p := &g.powerups[i]
		p.Draw()
	}

	for i := range g.sharks {
		s := &g.sharks[i]
		s.Draw()
	}

	for i := range g.pirates {
		p := &g.pirates[i]
		p.Draw()
	}

	for i := range g.seagulls {
		s := &g.seagulls[i]
		s.Draw()
	}

	if g.titanic != nil {
		g.titanic.Draw()
	}

	for i := range g.mines {
		m := &g.mines[i]
		m.Draw()
	}

	g.drawCannons(g.playerCannons)
	g.drawCannons(g.enemyCannons)

	g.player.Draw()
	g.ensemble.Draw()

	if g.paused {
		paused := "Paused"
		tw, th, _ := bigFont.SizeUTF8(paused)
		blitText(bigFont, (W-tw)/2, (H-th*4)/2, sdlcolor.Black, paused)
	}

	if g.gameOver != "" {
		tw, th, _ := bigFont.SizeUTF8(g.gameOver)
		blitText(bigFont, (W-tw)/2, (H-th)/2, sdlcolor.Black, g.gameOver)
	} else {
		message, _ := g.level.Message()
		if message != "" {
			tw, th, _ := smallFont.SizeUTF8(message)
			blitText(smallFont, (W-tw)/2, (H-th)/2, sdlcolor.Black, message)
		}
	}

	screen.Present()
}

func (g *Game) drawCannons(cannons []Cannon) {
	for i := range cannons {
		c := &cannons[i]
		c.Draw()
	}
}

func (g *Game) update() {
	if g.paused {
		return
	}

	if g.gameOver == "" {
		g.spawn()
	}

	g.updateEnemies()
	g.player.Update()
	g.health.Update()

	g.updateCannons(&g.playerCannons)
	g.updateCannons(&g.enemyCannons)

	g.score.Update()
	g.addEnvironmentEffects()
	g.ensemble.Update()
	g.updatePowerups()
	g.State.Update()

	if g.gameOver == "" {
		g.checkCollision()
	}

	switch {
	case g.health.Life == 0 && !g.player.Dying:
		g.player.Die()
	case g.player.Dying && !g.player.Dead:
		m := g.player.Image()
		p := Point{
			g.player.Pos.CenterX(m),
			g.player.Pos.CenterY(m),
		}

		g.ensemble.Explosion(p)
		g.ensemble.Debris(p)
	case g.player.Dead:
		g.gameOver = "Game Over"
	}

	g.lastShot++
	g.t++
}

func (g *Game) updateEnemies() {
	for i := 0; i < len(g.mines); {
		m := &g.mines[i]
		m.Update()

		remove := false
		if m.Exploding {
			if m.ExplodeFrames == 0 {
				remove = true
			}

			x := m.Pos.CenterX(m.Image())
			y := m.Pos.Y + float64(m.Image().W)/2
			p := Point{x, y}
			g.ensemble.Explosion(p)
			g.ensemble.Debris(p)
		}

		if m.Pos.Right(m.Image()) < 0 {
			remove = true
		}

		if remove {
			m.Free()
			l := len(g.mines) - 1
			g.mines[i], g.mines = g.mines[l], g.mines[:l]
		} else {
			i++
		}
	}

	for i := 0; i < len(g.sharks); {
		s := &g.sharks[i]
		s.Update()

		m := s.Image()
		if s.Dying {
			p := Point{s.Pos.CenterX(m), s.Pos.CenterY(m)}
			g.ensemble.Blood(p)
		}

		if s.Pos.Right(m) < 0 || s.Dead {
			s.Free()
			l := len(g.sharks) - 1
			g.sharks[i], g.sharks = g.sharks[l], g.sharks[:l]
		} else {
			i++
		}
	}

	for i := 0; i < len(g.pirates); {
		p := &g.pirates[i]
		p.Update()
		m := p.Image()

		center := Point{p.Pos.CenterX(m), p.Pos.CenterY(m)}
		if p.T%50 == 0 && !p.Dying {
			var c Cannon

			pos := Point{p.Pos.X, p.Pos.CenterY(m)}
			c.Init(pos, g.player.Angle, true, false)
			g.enemyCannons = append(g.enemyCannons, c)

			pt := p.Rotate(Point{0, 10}).Add(center)
			for i := 0; i < 4; i++ {
				g.ensemble.Fire(pt)
			}
		} else if p.Dying {
			g.ensemble.Explosion(center)
			g.ensemble.Wood(center)
		}

		if p.Pos.Right(m) < 0 || p.Dead {
			p.Free()
			l := len(g.pirates) - 1
			g.pirates[i], g.pirates = g.pirates[l], g.pirates[:l]
		} else {
			i++
		}
	}

	if g.titanic != nil {
		t := g.titanic
		t.Update()

		shoot := false
		angle := 0.0
		if !t.Dying {
			switch t.T % 100 {
			case 0:
				angle = 50
				shoot = true
			case 50:
				angle = 52.5
				shoot = true
			}
		}

		if shoot {
			for i := 0; i < 3; i++ {
				var c Cannon
				pos := Point{t.Pos.X, t.Pos.CenterY(t.Image())}
				c.Init(pos, g.player.Angle+float64(i-1)*10-angle, true, false)
				g.enemyCannons = append(g.enemyCannons, c)
			}
		}

		if t.Dead {
			g.gameOver = "Congratulations!\nYou sunk Titanic!"
			t.Free()
			g.titanic = nil
		}
	}

	for i := 0; i < len(g.seagulls); {
		s := &g.seagulls[i]
		s.Update()
		if s.Pos.Right(s.Image()) < 0 || s.Dead {
			s.Free()
			l := len(g.seagulls) - 1
			g.seagulls[i], g.seagulls = g.seagulls[l], g.seagulls[:l]
		} else {
			i++
		}
	}
}

func (g *Game) updateCannons(cannons *[]Cannon) {
	for i := 0; i < len(*cannons); {
		c := &(*cannons)[i]
		if !c.Underwater {
			p := c.Tail()
			q := p.Add(c.Vel.Scale(.5))
			g.ensemble.Trace(p)
			g.ensemble.Trace(q)
		}

		if c.Special && (!c.Underwater || rand.Float64() > 0.6) {
			p := c.Tail()
			g.ensemble.Explosion(p)
		}

		undOld := c.Underwater
		c.Update()
		if c.Underwater && !undOld {
			for i := 0; i < 5; i++ {
				p := Point{
					c.Pos.Right(c.Image()) - 4 + rand.Float64()*8,
					c.Pos.Y + rand.Float64()*2,
				}
				g.ensemble.Water(p)
			}
		}

		switch {
		case c.Pos.Right(c.Image()) < 0 && c.Vel.X < 0,
			c.Pos.X > W && c.Vel.X > 0,
			c.Vel.Y > H:
			l := len(*cannons) - 1
			(*cannons)[i], *cannons = (*cannons)[l], (*cannons)[:l]
		default:
			i++
		}
	}
}

func (g *Game) addEnvironmentEffects() {
	c := Point{
		g.player.Pos.CenterX(g.player.Image()),
		g.player.Pos.CenterY(g.player.Image()),
	}

	p := g.player.Rotate(Point{5 + rand.Float64()*9, 0})
	p = p.Add(c)

	q := g.player.Rotate(Point{19 + rand.Float64()*7, 5})
	q = q.Add(c)

	if !(g.spacePressed != 0 && g.t > g.spacePressed+Fps*3) {
		g.ensemble.Steam(p)
		g.ensemble.Steam(q)
	}

	if g.player.Dying && !g.player.Dead {
		g.ensemble.Explosion(c)
		g.ensemble.Debris(c)
	}

	if g.titanic != nil {
		for i := 0; i < 4; i++ {
			c := Point{
				g.titanic.Pos.CenterX(g.titanic.Image()),
				g.titanic.Pos.CenterY(g.titanic.Image()),
			}
			p := g.titanic.Rotate(Point{49 + rand.Float64()*9 + 28*float64(i), 25})
			p = p.Add(c)
			g.ensemble.Steam(p)
		}
	}

	if g.player.Splash {
		for i := 0; i < 10; i++ {
			r := rand.Float64()
			x := Lerp(g.player.Pos.X, g.player.Pos.Right(g.player.Image()), r)
			p := Point{x, WaterLevel(x)}
			g.ensemble.Water(p)
		}
	}
}

func (g *Game) updatePowerups() {
	for i := 0; i < len(g.powerups); {
		p := &g.powerups[i]
		p.Update()
		if p.Picked {
			l := len(g.powerups) - 1
			g.powerups[i], g.powerups = g.powerups[l], g.powerups[:l]
			p.Free()
		} else {
			i++
		}
	}
}

func (g *Game) event() {
	g.NextFrame = false
	for !g.NextFrame {
		select {
		case <-frame.C:
			g.NextFrame = true
		default:
		}

		for {
			ev := sdl.PollEvent()
			if ev == nil {
				break
			}
			switch ev := ev.(type) {
			case sdl.QuitEvent:
				g.Quit()

			case sdl.KeyDownEvent:
				switch ev.Sym {
				case sdl.K_ESCAPE:
					g.Quit()
				case sdl.K_p, sdl.K_RETURN:
					g.paused = !g.paused
				case sdl.K_s:
					Snapshot()
				case sdl.K_SPACE:
					g.playerFire()
				}

			case sdl.KeyUpEvent:
				switch ev.Sym {
				case sdl.K_LEFT:
					g.player.MoveLeft(false)
				case sdl.K_RIGHT:
					g.player.MoveRight(false)
				}
			}

			if !g.paused {
				switch ev := ev.(type) {
				case sdl.KeyDownEvent:
					switch ev.Sym {
					case sdl.K_LEFT:
						g.player.MoveLeft(true)
					case sdl.K_RIGHT:
						g.player.MoveRight(true)
					case sdl.K_UP:
						g.player.Jump()
					}
				}
			}
		}
	}
}

func (g *Game) spawn() {
	s := g.level.Spawn()

	if s&0x1 != 0 {
		s := Shark{}
		s.Init()
		g.sharks = append(g.sharks, s)
	}

	if s&0x2 != 0 {
		p := Pirate{}
		p.Init()
		g.pirates = append(g.pirates, p)
	}

	if s&0x4 != 0 {
		m := Mine{}
		m.Init()
		g.mines = append(g.mines, m)
	}

	if s&0x8 != 0 {
		s := Seagull{}
		s.Init()
		g.seagulls = append(g.seagulls, s)
	}

	if s&0x10 != 0 && g.titanic == nil {
		g.titanic = &Titanic{}
		g.titanic.Init()
	}

	if s&0x20 != 0 {
		p := Powerup{}
		p.Init()
		g.powerups = append(g.powerups, p)
	}
}

func (g *Game) playerFire() {
	if g.lastShot > MinFireDelay && !g.player.Dying {
		var c Cannon

		pos := Point{
			g.player.Pos.Right(g.player.Image()),
			g.player.Pos.CenterY(g.player.Image()),
		}
		c.Init(pos, g.player.Angle, false, false)
		g.playerCannons = append(g.playerCannons, c)

		g.lastShot = 0
		g.spacePressed = g.t
	}
}

func (g *Game) damagePlayer() {
	p := &g.player
	if !config.Invincibility {
		g.health.Damage()
		for i := 0; i < 10; i++ {
			pt := Point{rand.Float64() * 26, rand.Float64() * 10}
			ct := Point{p.Pos.CenterX(p.Image()), p.Pos.CenterY(p.Image())}
			pt = pt.Add(ct)
			g.ensemble.Debris(pt)
		}
	}
	p.Blinks += 12
}

func (g *Game) checkCollision() {
	p := &g.player

	for i := range g.powerups {
		pw := &g.powerups[i]
		if !pw.Fading && !p.Dying && Collision(&p.Entity, &pw.Entity) {
			g.health.Add()
			pw.Pickup()
		}
	}

	for i := range g.mines {
		m := &g.mines[i]
		if !m.Exploding && !p.Dying && Collision(&p.Entity, &m.Entity) {
			g.damagePlayer()
			m.Explode()
		}
	}

	for i := range g.sharks {
		s := &g.sharks[i]
		if !s.Dying && !p.Dying && Collision(&p.Entity, &s.Entity) {
			g.damagePlayer()
			s.Damage(1)
		}
	}

	g.checkCollisionPlayerCannon(&g.playerCannons)
	g.checkCollisionPlayerCannon(&g.enemyCannons)

	for i := range g.sharks {
		s := &g.sharks[i]
		for j := 0; j < len(g.playerCannons); {
			c := &g.playerCannons[j]
			if !s.Dying && Collision(&s.Entity, &c.Entity) {
				g.score.Add(15)
				s.Damage(1)
				s.Vel.X += c.Vel.X * 0.6
				s.Vel.Y += c.Vel.Y * 0.4
				if !c.Special {
					c.Free()
					l := len(g.playerCannons) - 1
					g.playerCannons[j], g.playerCannons = g.playerCannons[l], g.playerCannons[:l]
				}
				break
			} else {
				j++
			}
		}
	}

	for i := range g.seagulls {
		s := &g.seagulls[i]
		for j := 0; j < len(g.playerCannons); {
			c := &g.playerCannons[j]
			if !s.Dying && Collision(&s.Entity, &c.Entity) {
				g.score.Add(75)
				s.Damage(1)
				s.Vel.X += c.Vel.X * 0.6
				s.Vel.Y += c.Vel.Y * 0.4
				if !c.Special {
					c.Free()
					l := len(g.playerCannons) - 1
					g.playerCannons[j], g.playerCannons = g.playerCannons[l], g.playerCannons[:l]
				}
				break
			} else {
				j++
			}
		}
	}

	for i := range g.pirates {
		pe := &g.pirates[i]
		for j := 0; j < len(g.playerCannons); {
			c := &g.playerCannons[j]
			if !pe.Dying && Collision(&pe.Entity, &c.Entity) {
				snd := LoadSound("poks")
				snd.Play(0)

				g.score.Add(25)
				pe.Damage(1)

				for i := 0; i < 6; i++ {
					pt := Point{
						p.Pos.CenterX(p.Image()),
						p.Pos.CenterY(p.Image()),
					}
					pt.X += rand.Float64() * 15
					pt.Y += rand.Float64()*30 - 10
					g.ensemble.Wood(pt)
				}

				if !c.Special {
					c.Free()
					l := len(g.playerCannons) - 1
					g.playerCannons[j], g.playerCannons = g.playerCannons[l], g.playerCannons[:l]
				}
				break
			} else {
				j++
			}
		}
	}

	if g.titanic != nil {
		t := g.titanic
		for j := 0; j < len(g.playerCannons); {
			c := &g.playerCannons[j]
			if !t.Dying && Collision(&t.Entity, &c.Entity) {
				snd := LoadSound("poks")
				snd.Play(0)

				if c.Special {
					t.Damage(12)
					g.score.Add(100)
				} else {
					t.Damage(1)
					g.score.Add(7)
				}

				c.Free()
				l := len(g.playerCannons) - 1
				g.playerCannons[j], g.playerCannons = g.playerCannons[l], g.playerCannons[:l]
				break
			} else {
				j++
			}
		}
	}
}

func (g *Game) checkCollisionPlayerCannon(cannons *[]Cannon) {
	p := &g.player
	for i := 0; i < len(*cannons); {
		c := &(*cannons)[i]
		if !p.Dying && Collision(&p.Entity, &c.Entity) {
			g.damagePlayer()
			c.Free()
			l := len(*cannons) - 1
			(*cannons)[i], *cannons = (*cannons)[l], (*cannons)[:l]
		} else {
			i++
		}
	}
}

func (g *Game) clear() {
	g.player.Free()
	g.ensemble.Free()

	for i := range g.playerCannons {
		c := &g.playerCannons[i]
		c.Free()
	}

	for i := range g.enemyCannons {
		c := &g.enemyCannons[i]
		c.Free()
	}

	for i := range g.pirates {
		p := &g.pirates[i]
		p.Free()
	}

	for i := range g.mines {
		m := &g.mines[i]
		m.Free()
	}

	for i := range g.seagulls {
		s := &g.seagulls[i]
		s.Free()
	}

	for i := range g.sharks {
		s := &g.sharks[i]
		s.Free()
	}

	for i := range g.powerups {
		p := &g.powerups[i]
		p.Free()
	}

	if g.titanic != nil {
		g.titanic.Free()
		g.titanic = nil
	}

	g.playerCannons = g.playerCannons[:0]
	g.enemyCannons = g.enemyCannons[:0]
	g.pirates = g.pirates[:0]
	g.mines = g.mines[:0]
	g.seagulls = g.seagulls[:0]
	g.sharks = g.sharks[:0]
	g.powerups = g.powerups[:0]
}
