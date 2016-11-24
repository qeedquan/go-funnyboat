package main

import (
	"bytes"
	"math/rand"
	"text/template"
)

type Map struct {
	Length  []int
	Message []string
	Color   []string
	Weather []float64
	Phase   [][][]int
}

var normalMap = Map{
	Length: []int{900, 900, 900, 1800, 1300, -1},

	Message: []string{
		"Watch out for those angry sharks, captain!",
		"Minefield ahead, captain!",
		"Oh no! It's the infamous fleet of pirate Captain {{.Color}}beard!",
		"",
		"Uh, oh. Looks like some busy waters ahead, captain!",
		"Holy cow! It's the legendary Titanic!",
	},

	Color: []string{"Brown", "Red", "Yellow", "Magenta", "Pink", "Cyan",
		"Blue", "Black", "Green", "Violet", "Beige", "White", "Gray",
		"Blonde", "Orange", "Brunette", "Ginger", "Turquoise"},

	Weather: []float64{30, 40, 20, 30, 60, 5},

	Phase: [][][]int{
		{
			{10, 80},
			{230, 0},
			{120, 0},
			{0, 3000},
			{1, 0},
			{450, 1000},
		},
		{
			{100, 300},
			{257, 0},
			{70, 137},
			{0, 3000},
			{1, 0},
			{0, 1000},
		},
		{
			{257, 500},
			{30, 300},
			{470, 500},
			{0, 1500},
			{1, 0},
			{0, 1000},
		},
		{
			{0, 183},
			{230, 319},
			{40, 217},
			{0, 700},
			{1, 0},
			{0, 1000},
		},
		{
			{0, 233},
			{230, 519},
			{40, 317},
			{0, 700},
			{1, 0},
			{0, 1000},
		},
		{
			{70, 200},
			{300, 0},
			{0, 200},
			{0, 0},
			{10, -1},
			{0, 0},
		},
	},
}

var endlessMap = Map{
	Length: []int{450},

	Message: []string{"This is the endless mode.\nGood luck!"},

	Weather: []float64{30, 10, 50},

	Phase: [][][]int{
		{
			{0, 255},
			{150, 257},
			{50, 253},
			{0, 507},
			{0, 0},
			{100, 0},
		},

		{
			{0, 150},
			{400, 700},
			{50, 700},
			{0, 507},
			{0, 0},
			{500, 0},
		},

		{
			{150, 400},
			{0, 150},
			{350, 500},
			{0, 507},
			{0, 0},
			{500, 0},
		},

		{
			{350, 500},
			{150, 400},
			{20, 150},
			{0, 507},
			{0, 0},
			{100, -1},
		},
	},
}

type Level struct {
	endless bool
	text    string
	phase   int
	t       int
}

func (l *Level) Reset(endless bool) {
	l.endless = endless
	l.phase = 0
	l.text = ""
	l.t = 0
}

func (l *Level) Spawn() uint8 {
	var s uint8

	m := l.curmap()
	ln := m.Length[l.phase%len(m.Length)]
	if l.endless && ln < 30 {
		ln = 30
	}

	if ln != -1 && l.t > ln {
		l.t = 0
		mod := len(m.Length) * len(m.Weather) * len(m.Phase)
		l.phase = (l.phase + 1) % mod
	}

	i := uint(0)
	for _, enemy := range m.Phase[l.phase%len(m.Phase)] {
		offset, delay := enemy[0], enemy[1]
		if l.endless && delay > 0 {
			delay -= l.phase / 4 * 5
			offset -= l.phase / 4 * 5
			if delay <= 30 {
				delay = 30
			}
			if offset < 0 {
				offset = 0
			}
		}

		if (delay == -1 && l.t == offset) || (delay != 0 && l.t%delay == offset) {
			s |= (1 << i)
		}
		i++
	}

	w := l.phase
	if l.endless {
		w /= 4
	}
	SetWaterAmplitude(m.Weather[w%len(m.Weather)])

	l.t++

	return s
}

func (l *Level) Message() (text string, alpha uint8) {
	text = l.text
	alpha = 0

	m := l.curmap()
	if l.t < 120 && l.phase < len(m.Message) {
		alpha = 255
		if l.t > 60 {
			alpha = uint8(255 - (l.t-60)*255/60)
		}

		if text != "" {
			return
		}

		buf := new(bytes.Buffer)
		message := m.Message[l.phase]
		tmpl := template.Must(template.New("").Parse(message))
		err := tmpl.Execute(buf, l)
		if err != nil {
			panic(err)
		}
		l.text = buf.String()
		text = l.text

		return
	}

	l.text = ""
	return
}

func (l *Level) curmap() *Map {
	if l.endless {
		return &endlessMap
	}
	return &normalMap
}

func (l *Level) Color() string {
	m := l.curmap()
	if len(m.Color) == 0 {
		return "Color"
	}
	return m.Color[rand.Intn(len(m.Color))]
}
