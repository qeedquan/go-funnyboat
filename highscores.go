package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/qeedquan/go-media/sdl"
	"github.com/qeedquan/go-media/sdl/sdlimage/sdlcolor"
)

type Rank struct {
	Name  string
	Value int
}

type RankSlice []Rank

func (s RankSlice) Len() int           { return len(s) }
func (s RankSlice) Less(i, j int) bool { return s[i].Value > s[j].Value }
func (s RankSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

type Highscores struct {
	State
	ranks    []Rank
	title    string
	filename string
}

var dummyScores = []Rank{
	{"Funny Boat", 2000},
	{"Hectigo", 1500},
	{"JDruid", 1000},
	{"Pekuja", 750},
	{"Pirate", 500},
	{"Shark", 400},
	{"Seagull", 300},
	{"Naval Mine", 200},
	{"Cannonball", 100},
	{"Puffy the Cloud", 50},
}

func (h *Highscores) Reset(endless bool, newScore int) {
	h.State.Reset()

	h.title = "Story Mode"
	h.filename = "scores"
	if endless {
		h.title = "Endless Mode"
		h.filename = "endless_scores"
	}

	h.Load()
	h.Update(newScore)
}

func (h *Highscores) Load() {
	var filename string
	var err error

	h.ranks = h.ranks[:0]

	log.SetPrefix("scores: ")
	defer func() {
		if err != nil {
			log.Print("load failure: ", err)
			h.ranks = append(h.ranks[:0], dummyScores...)
		} else {
			log.Printf("load %q", filename)
		}

		if len(h.ranks) >= MaxRanks {
			h.ranks = h.ranks[:MaxRanks]
		}
	}()

	path, err := config.Path()
	if err != nil {
		return
	}

	filename = filepath.Join(path, h.filename)
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	var line [2]string

	s := bufio.NewScanner(f)
loop:
	for n := 0; n < MaxRanks; n++ {
		for i := 0; i < 2; i++ {
			if !s.Scan() {
				break loop
			}
			line[i] = s.Text()
		}

		value, err := strconv.Atoi(line[1])
		if err != nil {
			continue
		}
		name := strings.TrimSpace(line[0])
		h.ranks = append(h.ranks, Rank{name, value})
	}

	sort.Stable(RankSlice(h.ranks))
}

func (h *Highscores) Save() {
	var filename string
	var err error

	log.SetPrefix("scores: ")
	defer func() {
		if err != nil {
			log.Print("save error: ", err)
		} else {
			log.Printf("saved to %q", filename)
		}
	}()

	path, err := config.Path()
	if err != nil {
		return
	}

	filename = filepath.Join(path, h.filename)
	f, err := os.Create(filename)
	if err != nil {
		return
	}

	w := bufio.NewWriter(f)
	for _, r := range h.ranks {
		fmt.Fprintf(w, "%v\n%v\n", r.Name, r.Value)
	}
	flushErr := w.Flush()
	closeErr := f.Close()

	err = flushErr
	if err == nil {
		err = closeErr
	}
}

func (h *Highscores) Update(newScore int) {
	if newScore < 0 || (len(h.ranks) > 0 && newScore < h.ranks[len(h.ranks)-1].Value) {
		return
	}
	h.ranks = append(h.ranks, Rank{config.Name, newScore})
	sort.Stable(RankSlice(h.ranks))
	if len(h.ranks) >= MaxRanks {
		h.ranks = h.ranks[:MaxRanks]
	}
	h.Save()
}

func (h *Highscores) Run(endless bool, newScore int) {
	h.Reset(endless, newScore)
	for !h.Done {
		h.draw()
		h.State.Update()
		h.event()
	}
}

func (h *Highscores) draw() {
	h.State.Draw()

	tw, _, _ := bigFont.SizeUTF8(h.title)
	blitText(bigFont, (W-tw)/2, 10, sdlcolor.Black, fmt.Sprint(h.title))

	for i, r := range h.ranks {
		_, th, _ := smallFont.SizeUTF8(r.Name)
		x := 10
		y := 50 + i*th
		blitText(smallFont, x, y, sdlcolor.Black, fmt.Sprintf("%v. %v", i+1, r.Name))

		value := fmt.Sprint(r.Value)
		tw, _, _ := smallFont.SizeUTF8(value)
		x = W - tw - 10
		blitText(smallFont, x, y, sdlcolor.Black, fmt.Sprint(value))
	}

	screen.Present()
}

func (h *Highscores) event() {
	h.NextFrame = false
	for !h.NextFrame {
		select {
		case <-frame.C:
			h.NextFrame = true
		default:
		}

		for {
			ev := sdl.PollEvent()
			if ev == nil {
				break
			}
			switch ev.(type) {
			case sdl.QuitEvent, sdl.KeyDownEvent:
				h.Save()
				h.Quit()
			}
		}
	}
}
