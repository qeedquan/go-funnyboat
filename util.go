package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"path/filepath"

	"github.com/qeedquan/go-media/sdl"
)

const (
	Degree = 180 / math.Pi
	Radian = math.Pi / 180
)

func Lerp(a, b, t float64) float64 {
	return a*t + b*(1-t)
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Snapshot() {
	var err error
	var filename string

	log.SetPrefix("snapshot: ")
	defer func() {
		if err != nil {
			log.Print("failure: ", err)
		} else {
			log.Printf("saved to %q", filename)
		}
	}()

	path, err := config.Path()
	if err != nil {
		return
	}

	path = filepath.Join(path, "ss")
	err = os.MkdirAll(path, 0755)
	if err != nil && !os.IsExist(err) {
		return
	}

	glob := filepath.Join(path, "ss_*.png")
	matches, err := filepath.Glob(glob)
	if err != nil {
		return
	}

	var v uint64
	i := uint64(0)
	for _, m := range matches {
		b := filepath.Base(m)
		n, _ := fmt.Sscanf(b, "ss_%v", &v)
		if n == 1 {
			if i <= v {
				i = v + 1
			}
		}
	}

	filename = filepath.Join(path, fmt.Sprint("ss_", i, ".png"))
	f, err := os.Create(filename)
	if err != nil {
		return
	}
	defer func() {
		closeErr := f.Close()
		if err == nil {
			err = closeErr
		}
	}()

	w, h, err := screen.OutputSize()
	if err != nil {
		return
	}

	stride := 4 * w
	pixels := make([]byte, stride*h)
	err = screen.ReadPixels(nil, sdl.PIXELFORMAT_ABGR8888, pixels, stride)
	if err != nil {
		return
	}

	rgba := &image.RGBA{
		Stride: stride,
		Pix:    pixels,
		Rect:   image.Rect(0, 0, w, h),
	}

	r := rgba.Bounds()
	for x := r.Max.X; x >= r.Min.X; x-- {
		if (rgba.RGBAAt(x, 0) != color.RGBA{}) {
			rgba = rgba.SubImage(image.Rect(r.Min.X, r.Min.Y, x, r.Max.Y)).(*image.RGBA)
			break
		}
	}

	err = png.Encode(f, rgba)
	if err != nil {
		return
	}
}
