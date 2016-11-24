package main

import (
	"image"
	"image/color"
	"image/draw"
	"log"
	"math"
	"os"
	"path/filepath"

	"github.com/qeedquan/go-media/sdl"
	"github.com/qeedquan/go-media/sdl/sdlimage"
	"github.com/qeedquan/go-media/sdl/sdlimage/sdlcolor"
	"github.com/qeedquan/go-media/sdl/sdlttf"
)

type Point struct {
	X, Y float64
}

func (p Point) Add(q Point) Point {
	return Point{p.X + q.X, p.Y + q.Y}
}

func (p Point) CenterX(i *Image) float64 {
	return (p.X + p.X + float64(i.W)) / 2
}

func (p Point) CenterY(i *Image) float64 {
	return (p.Y + p.Y + float64(i.H)) / 2
}

func (p Point) Right(i *Image) float64 {
	return p.X + float64(i.W)
}

func (p Point) Bottom(i *Image) float64 {
	return p.Y + float64(i.H)
}

func (p Point) Scale(s float64) Point {
	return Point{p.X * s, p.Y * s}
}

type Image struct {
	*sdl.Texture
	Store  *image.Alpha
	Buffer *image.Alpha
	Alpha  *image.Alpha
	Angle  float64
	W, H   int
}

func (i *Image) Blit(pos Point) {
	screen.CopyEx(i.Texture, nil, &sdl.Rect{int32(pos.X), int32(pos.Y), int32(i.W), int32(i.H)}, -i.Angle, nil, sdl.FLIP_NONE)
}

func (i *Image) Copy() *Image {
	return i.CopySize(i.W, i.H)
}

func (i *Image) CopySize(w, h int) *Image {
	m := NewImage(w, h)
	m.Bind()
	i.Blit(Point{})
	m.Unbind()
	draw.Draw(m.Store, i.Store.Bounds(), i.Store, image.ZP, draw.Src)
	draw.Draw(m.Alpha, i.Alpha.Bounds(), i.Alpha, image.ZP, draw.Src)
	return m
}

func (i *Image) Bind() {
	log.SetPrefix("image: ")
	err := screen.SetTarget(i.Texture)
	if err != nil {
		log.Fatal(err)
	}
}

func (i *Image) Unbind() {
	log.SetPrefix("image: ")
	err := screen.SetTarget(nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (i *Image) Vline(x, y1, y2 int) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}

	screen.DrawLine(x, y1, x, y2)
	for y := y1; y <= y2; y++ {
		i.Store.SetAlpha(x, y, color.Alpha{255})
		i.Alpha.SetAlpha(x, y, color.Alpha{255})
	}
}

func (i *Image) UpdateAngle(angle float64) {
	i.Angle = angle
	draw.Draw(i.Buffer, i.Buffer.Bounds(), &image.Uniform{color.Transparent}, image.ZP, draw.Src)
	i.Alpha = rotateAlpha(i.Store, i.Buffer, i.Angle)
}

func LoadSurface(name string) (*sdl.Surface, error) {
	filename := filepath.Join(config.Resource, name+".png")
	return sdlimage.LoadSurfaceFile(filename)
}

var (
	images = make(map[string]*Image)
)

func NewImage(w, h int) *Image {
	log.SetPrefix("image: ")

	texture, err := screen.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_TARGET, w, h)
	if err != nil {
		log.Fatal(err)
	}
	texture.SetBlendMode(sdl.BLENDMODE_BLEND)

	max := Max(w, h) * 2
	store := image.NewAlpha(image.Rect(0, 0, w, h))
	buffer := image.NewAlpha(image.Rect(0, 0, max, max))
	alpha := buffer.SubImage(image.Rect(0, 0, w, h)).(*image.Alpha)

	m := &Image{
		Texture: texture,
		Store:   store,
		Buffer:  buffer,
		Alpha:   alpha,
		W:       w,
		H:       h,
	}

	m.Bind()
	screen.SetDrawColor(sdlcolor.Transparent)
	screen.Clear()
	m.Unbind()

	return m
}

func LoadImage(name string) *Image {
	if m, found := images[name]; found {
		return m
	}

	log.SetPrefix("image: ")
	filename := filepath.Join(config.Resource, name+".png")

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	texture, err := sdlimage.LoadTextureImage(screen.Renderer, img)
	if err != nil {
		log.Fatal(err)
	}
	texture.SetBlendMode(sdl.BLENDMODE_BLEND)

	_, _, w, h, _ := texture.Query()

	max := Max(w, h) * 2
	buffer := image.NewAlpha(image.Rect(0, 0, max, max))
	store := image.NewAlpha(image.Rect(0, 0, w, h))
	alpha := buffer.SubImage(image.Rect(0, 0, w, h)).(*image.Alpha)
	draw.Draw(store, img.Bounds(), img, image.ZP, draw.Src)
	draw.Draw(alpha, img.Bounds(), img, image.ZP, draw.Src)

	m := &Image{
		Texture: texture,
		Store:   store,
		Buffer:  buffer,
		Alpha:   alpha,
		W:       w,
		H:       h,
	}
	images[name] = m
	return m
}

func LoadFont(name string, ptsize int) *sdlttf.Font {
	log.SetPrefix("font: ")
	filename := filepath.Join(config.Resource, name+".ttf")
	f, err := sdlttf.OpenFont(filename, ptsize)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func rotate90(src *image.Alpha, buf *image.Alpha, angle int) (dst *image.Alpha) {
	const bpp = 1

	nturns := (angle / 90) % 4
	if nturns < 0 {
		nturns += 4
	}

	sb := src.Bounds()
	if nturns%2 == 0 {
		dst = buf.SubImage(image.Rect(0, 0, sb.Dx(), sb.Dy())).(*image.Alpha)
	} else {
		dst = buf.SubImage(image.Rect(0, 0, sb.Dy(), sb.Dx())).(*image.Alpha)
	}
	db := dst.Bounds()

	sx, dx := bpp, bpp
	sy, dy := src.Stride, dst.Stride
	sp, dp := 0, 0

	switch nturns {
	case 0: // identity
	case 1:
		sp = (sb.Dx() - 1) * sx
		sy = -sx
		sx = src.Stride
	case 2:
		sp = (sb.Dy()-1)*sy + (sb.Dx()-1)*sx
		sx = -sx
		sy = -sy
	case 3:
		sp = (sb.Dy() - 1) * sy
		sx = -sy
		sy = bpp
	}

	for ly := db.Min.Y; ly < db.Max.Y; ly++ {
		i, j := dp, sp
		for lx := db.Min.X; lx < db.Max.X; lx++ {
			copy(dst.Pix[i:i+bpp], src.Pix[j:j+bpp])
			i += dx
			j += sx
		}
		dp += dy
		sp += sy
	}

	return
}

func newRotationImage(src *image.Alpha, buf *image.Alpha, angle float64) (dst *image.Alpha, sin, cos float64) {
	rad := angle * math.Pi / 180
	sin, cos = math.Sincos(rad)

	sb := src.Bounds()
	cx := cos * float64(sb.Dx())
	cy := cos * float64(sb.Dy())
	sx := sin * float64(sb.Dx())
	sy := sin * float64(sb.Dy())

	a1 := math.Abs(cx + sy)
	a2 := math.Abs(cx - sy)
	a3 := math.Abs(-cx + sy)
	a4 := math.Abs(-cx - sy)
	width := math.Max(a1, a2)
	width = math.Max(width, a3)
	width = math.Max(width, a4)

	a1 = math.Abs(sx + cy)
	a2 = math.Abs(sx - cy)
	a3 = math.Abs(-sx + cy)
	a4 = math.Abs(-sx - cy)
	height := math.Max(a1, a2)
	height = math.Max(width, a3)
	height = math.Max(width, a4)

	dst = buf.SubImage(image.Rect(0, 0, int(width), int(height))).(*image.Alpha)
	return
}

func rotateAlpha(src *image.Alpha, buf *image.Alpha, angle float64) (dst *image.Alpha) {
	if integral, frac := math.Modf(angle); int(integral) == 90 && frac == 0 {
		return rotate90(src, buf, int(integral))
	}

	const bpp = 1
	dst, sin, cos := newRotationImage(src, buf, angle)
	isin, icos := int(sin*65536), int(cos*65536)

	sb, db := src.Bounds(), dst.Bounds()
	sw, sh, dw, dh := sb.Dx(), sb.Dy(), db.Dx(), db.Dy()

	cy := dh / 2
	xd := (sw - dw) << 15
	yd := (sh - dh) << 15

	ax := (dw << 15) - int(cos*float64((dw-1)<<15))
	ay := (dh << 15) - int(sin*float64((dw-1)<<15))

	xmax := (sw << 16) - 1
	ymax := (sh << 16) - 1

	dp := 0
	s, d := src.Pix, dst.Pix
	for ly := db.Min.Y; ly < db.Max.Y; ly++ {
		dx := (ax + (isin * (cy - ly - db.Min.Y))) + xd
		dy := (ay - (icos * (cy - ly - db.Min.Y))) + yd
		i := dp
		for lx := db.Min.X; lx < db.Max.X; lx++ {
			if dx < 0 || dy < 0 || dx > xmax || dy > ymax {
				for j := 0; j < bpp; j++ {
					d[i+j] = 0
				}
			} else {
				sy := (dy >> 16) * src.Stride
				sx := (dx >> 16) * bpp
				for j := 0; j < bpp; j++ {
					d[i+j] = s[sy+sx+j]
				}
			}

			i += bpp
			dx += icos
			dy += isin
		}
		dp += dst.Stride
	}
	return
}

func blitText(font *sdlttf.Font, x, y int, c sdl.Color, text string) {
	log.SetPrefix("text: ")
	r, err := font.RenderUTF8BlendedEx(surface, text, c)
	if err != nil {
		log.Fatal(err)
	}

	p, err := texture.Lock(nil)
	if err != nil {
		log.Fatal(err)
	}

	err = surface.Lock()
	if err != nil {
		log.Fatal(err)
	}
	s := surface.Pixels()
	for i := 0; i < len(p); i += 4 {
		p[i] = s[i+2]
		p[i+1] = s[i]
		p[i+2] = s[i+1]
		p[i+3] = s[i+3]
	}

	surface.Unlock()
	texture.Unlock()

	texture.SetBlendMode(sdl.BLENDMODE_BLEND)
	screen.Copy(texture, &sdl.Rect{0, 0, r.W, r.H}, &sdl.Rect{int32(x), int32(y), r.W, r.H})
}