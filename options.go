package main

import (
	"fmt"
	"unicode/utf8"

	"github.com/qeedquan/go-media/sdl"
)

type toggle bool

func (t toggle) String() string {
	if t {
		return "on"
	}
	return "off"
}

type Options struct {
	Selector
}

func (o *Options) Init() {
	o.Selector.Init(smallFont, o.event)
}

func (o *Options) Run() {
	o.refresh()
	o.Selector.Run(o.menu, 0)
	config.Save()
}

func (o *Options) event(ev sdl.Event) {
	switch ev := ev.(type) {
	case sdl.KeyDownEvent:
		switch o.cursor {
		case 3:
			switch ev.Sym {
			case sdl.K_BACKSPACE:
				o.updateName("", true)
			case sdl.K_SPACE:
				o.updateName(" ", false)
			}
		default:
			switch ev.Sym {
			case sdl.K_RETURN, sdl.K_KP_ENTER, sdl.K_SPACE:
				o.toggle()
			}
		}
	case sdl.TextInputEvent:
		if o.cursor == 3 {
			for i := range ev.Text {
				if ev.Text[i] == 0 {
					o.updateName(string(ev.Text[:i]), false)
					break
				}
			}
		}
	}
	o.refresh()
}

func (o *Options) refresh() {
	c := &config
	o.menu = []string{
		fmt.Sprint("Particle effects: ", toggle(c.Particles)),
		fmt.Sprint("Sound effects: ", toggle(c.Sound)),
		fmt.Sprint("Music: ", toggle(c.Music)),
		fmt.Sprint("Player Name: ", c.Name),
		fmt.Sprint("Invincibility: ", toggle(c.Invincibility)),
	}
}

func (o *Options) toggle() {
	c := &config
	switch o.cursor {
	case 0:
		c.Particles = !c.Particles
	case 1:
		c.Sound = !c.Sound
	case 2:
		c.Music = !c.Music
		song.Play()
	case 4:
		c.Invincibility = !c.Invincibility
	}
}

func (o *Options) updateName(s string, backspace bool) {
	c := &config
	l := utf8.RuneCountInString(c.Name)
	switch {
	case backspace && l > 0:
		c.Name = c.Name[:len(c.Name)-1]
	case !backspace && l < MaxName:
		c.Name += s
	}
}
