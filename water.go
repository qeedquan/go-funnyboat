package main

import (
	"math"

	"github.com/qeedquan/go-media/sdl"
)

type Water struct {
	image   *Image
	levels  []float64
	ta, a   float64
	tw, w   float64
	ts, s   float64
	tbh, bh float64
	xm, tm  float64
	t       float64
}

var (
	water Water
)

func InitWater() {
	water.Init()
}

func UpdateWater() {
	water.Update()
}

func DrawWater() {
	water.Draw()
}

func SetWaterAmplitude(a float64) {
	water.SetAmplitude(a)
}

func WaterLevel(x float64) float64 {
	return water.Level(x)
}

func WaterRegion(p Point, i *Image) [3]float64 {
	return [3]float64{
		WaterLevel(p.X),
		WaterLevel(p.CenterX(i)),
		WaterLevel(p.Right(i)),
	}
}

func (w *Water) Init() {
	w.image = NewImage(W, H)
	w.levels = make([]float64, W)

	w.ta = H / 8
	w.tw = 0.02 * W / (2 * math.Pi)
	w.ts = 0.06 / (2 * math.Pi) * Fps
	w.tbh = H / 24 * 8
	w.a, w.w, w.s, w.bh = w.ta, w.tw, w.ts, w.tbh

	w.xm = 2 * math.Pi / w.w / W
	w.tm = 2 * math.Pi / Fps * w.s

	w.Update()
}

func (w *Water) Update() {
	w.image.Bind()
	defer w.image.Unbind()

	screen.SetDrawColor(sdl.Color{200, 210, 255, 0})
	screen.Clear()

	screen.SetDrawColor(sdl.Color{20, 60, 180, 110})
	for x := range w.levels {
		h := H - (math.Sin(float64(x)*w.xm+w.t*w.tm)*w.a + w.bh)
		w.levels[x] = h
		hi, _ := math.Modf(h)
		w.image.Vline(x, int(hi), H)
	}

	if w.ta != w.a {
		w.a = Lerp(w.a, w.ta, 0.99)
	}
	if w.tw != w.w {
		w.w = Lerp(w.w, w.tw, 0.99)
		w.xm = 2 * math.Pi / w.w / W
	}
	if w.ts != w.s {
		w.s = Lerp(w.s, w.ts, 0.99)
		w.tm = 2 * math.Pi / Fps * w.s
	}
	if w.tbh != w.bh {
		w.bh = Lerp(w.bh, w.tbh, 0.99)
	}

	w.t++
}

func (w *Water) Draw() {
	w.image.Blit(Point{})
}

func (w *Water) Level(x float64) float64 {
	xi := int(x)
	if xi >= len(w.levels) {
		xi = len(w.levels) - 1
	} else if xi < 0 {
		xi = 0
	}
	return w.levels[xi]
}

func (w *Water) SetAmplitude(a float64) {
	w.ta = a
}
