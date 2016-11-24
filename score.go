package main

import (
	"fmt"

	"github.com/qeedquan/go-media/sdl/sdlimage/sdlcolor"
)

type Score struct {
	target int
	Value  int
	pos    Point
}

func (s *Score) Reset() {
	s.target = 0
	s.Value = 0
	s.pos = Point{100, 5}
}

func (s *Score) Draw() {
	blitText(smallFont, int(s.pos.X), int(s.pos.Y), sdlcolor.Black, fmt.Sprintf("Score: %v", s.Value))
}

func (s *Score) Update() {
	if s.target > s.Value {
		s.Value++
	}
}

func (s *Score) Add(points int) {
	s.target += points
	s.Update()
}
