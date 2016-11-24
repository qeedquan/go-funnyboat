package main

import "github.com/qeedquan/go-media/sdl"

type Menu struct {
	Selector
}

func (m *Menu) Init() {
	m.Selector.Init(bigFont, m.event)
}

func (m *Menu) event(ev sdl.Event) {
	switch ev := ev.(type) {
	case sdl.KeyDownEvent:
		switch ev.Sym {
		case sdl.K_SPACE, sdl.K_RETURN, sdl.K_KP_ENTER:
			m.Quit()
		}
	}
}
