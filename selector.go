package main

import (
	"fmt"

	"github.com/qeedquan/go-media/sdl"
	"github.com/qeedquan/go-media/sdl/sdlimage/sdlcolor"
	"github.com/qeedquan/go-media/sdl/sdlttf"
)

type Selector struct {
	State
	logo    *Image
	font    *sdlttf.Font
	menu    []string
	cursor  int
	handler func(sdl.Event)
}

func (s *Selector) Init(font *sdlttf.Font, handler func(sdl.Event)) {
	s.State.Init()
	s.logo = LoadImage("logo")
	s.font = font
	s.handler = handler
}

func (s *Selector) Run(menu []string, cursor int) int {
	s.Reset()
	s.menu = menu
	s.cursor = cursor

	for !s.Done {
		s.draw()
		s.State.Update()
		s.event()
	}
	return s.cursor
}

func (s *Selector) draw() {
	s.State.Draw()

	for i := range s.menu {
		s.render(i)
	}

	p := Point{(W - float64(s.logo.W)) / 2, 0}
	s.logo.Blit(p)

	const link = "http://funnyboat.sourceforge.net/"
	tw, th, _ := smallFont.SizeUTF8(link)
	p = Point{(W - float64(tw)) / 2, H - float64(th)}
	blitText(smallFont, int(p.X), int(p.Y), sdlcolor.Black, link)

	screen.Present()
}

func (s *Selector) render(id int) {
	color := sdlcolor.Black
	if s.cursor == id {
		color = sdl.Color{255, 127, 0x00, 0xff}
	}

	title := s.menu[id]
	tw, th, _ := s.font.SizeUTF8(title)
	x := (W - tw) / 2
	y := s.logo.H + id*th

	blitText(s.font, x, y, color, fmt.Sprintf("%v", title))
}

func (s *Selector) event() {
	s.NextFrame = false
	for !s.NextFrame {
		select {
		case <-frame.C:
			s.NextFrame = true
		default:
		}

		for {
			ev := sdl.PollEvent()
			if ev == nil {
				break
			}
			switch ev := ev.(type) {
			case sdl.QuitEvent:
				s.quit()
			case sdl.KeyDownEvent:
				switch ev.Sym {
				case sdl.K_ESCAPE:
					s.quit()
				case sdl.K_DOWN:
					s.move(1)
				case sdl.K_UP:
					s.move(-1)
				default:
					s.handler(ev)
				}
			default:
				s.handler(ev)
			}
		}
	}
}

func (s *Selector) move(i int) {
	s.cursor = (s.cursor + i) % len(s.menu)
	if s.cursor < 0 {
		s.cursor += len(s.menu)
	}
}

func (s *Selector) quit() {
	s.cursor = -1
	s.Quit()
}
