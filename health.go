package main

type Health struct {
	heart    *Image
	empty    *Image
	broken   *Image
	Life     int
	lost     int
	counters [MaxHearts]int
}

func (h *Health) Init() {
	h.heart = LoadImage("sydan")
	h.empty = LoadImage("sydan-tyhja")
	h.broken = LoadImage("sydan-rikki")
}

func (h *Health) Reset() {
	h.Life = MaxHearts
	h.lost = 0
}

func (h *Health) Draw() {
	for i := 0; i < MaxHearts; i++ {
		p := Point{10 + float64(i*(h.heart.W+1)), float64(h.heart.H)}
		if i < h.Life {
			h.heart.Blit(p)
		} else {
			h.empty.Blit(p)
			if i < h.Life+h.lost {
				p.X -= float64(h.counters[i])
				h.broken.Blit(p)
			}
		}
	}
}

func (h *Health) Update() {
	for i := h.Life; i < MaxHearts; i++ {
		if i < h.Life+h.lost {
			if h.counters[i]++; h.counters[i] == 25 {
				h.lost--
			}
		}
	}
}

func (h *Health) Damage() {
	if h.Life > 0 {
		h.Life--
		h.lost++
		h.Update()
	}
}

func (h *Health) Add() {
	if h.Life < MaxHearts {
		h.Life++
	}
}
