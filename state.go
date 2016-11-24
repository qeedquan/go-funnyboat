package main

type State struct {
	Done      bool
	NextFrame bool
	Sky       *Image
}

func (s *State) Init() {
	s.Sky = LoadImage("taivas")
}

func (s *State) Reset() {
	s.Done, s.NextFrame = false, false
}

func (s *State) Quit() {
	s.Done, s.NextFrame = true, true
}

func (s *State) Update() {
	UpdateClouds()
	UpdateWater()
}

func (s *State) Draw() {
	s.Sky.Blit(Point{})
	DrawClouds()
	DrawWater()
}
