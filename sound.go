package main

import (
	"log"
	"path/filepath"

	"github.com/qeedquan/go-media/sdl/sdlmixer"
)

type Music struct {
	*sdlmixer.Music
}

type Sound struct {
	*sdlmixer.Chunk
}

func LoadMusic(name string) *Music {
	log.SetPrefix("music: ")

	filename := filepath.Join(config.Resource, name+".ogg")
	mus, err := sdlmixer.LoadMUS(filename)
	if err != nil {
		log.Print(err)
		return nil
	}

	return &Music{mus}
}

var (
	sounds = make(map[string]*Sound)
)

func LoadSound(name string) *Sound {
	if s, found := sounds[name]; found {
		return s
	}

	log.SetPrefix("sound: ")
	filename := filepath.Join(config.Resource, name+".ogg")
	chunk, err := sdlmixer.LoadWAV(filename)
	if err != nil {
		log.Print(err)
		return nil
	}

	s := &Sound{chunk}
	sounds[name] = s
	return s
}

func (m *Music) Play() {
	if m == nil || !config.Music {
		sdlmixer.HaltMusic()
		return
	}

	m.Music.Play(-1)
}

func (m *Music) Free() {
	if m == nil {
		return
	}
	m.Music.Free()
}

func (s *Sound) Play(loops int) {
	if s == nil || !config.Sound {
		return
	}
	s.PlayChannel(-1, loops)
}

func (s *Sound) Free() {}
