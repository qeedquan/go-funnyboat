package main

import (
	"log"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/qeedquan/go-media/sdl"
	"github.com/qeedquan/go-media/sdl/sdlimage/sdlcolor"
	"github.com/qeedquan/go-media/sdl/sdlmixer"
	"github.com/qeedquan/go-media/sdl/sdlttf"
)

type Display struct {
	*sdl.Window
	*sdl.Renderer
}

func newDisplay(w, h int, wflag sdl.WindowFlags) (*Display, error) {
	window, renderer, err := sdl.CreateWindowAndRenderer(w, h, wflag)
	if err != nil {
		return nil, err
	}
	return &Display{window, renderer}, nil
}

const (
	W            = 400
	H            = 300
	MaxHearts    = 5
	MaxName      = 32
	MaxRanks     = 10
	MinFireDelay = 1
	Fps          = 30
)

var (
	config Config

	screen    *Display
	smallFont *sdlttf.Font
	bigFont   *sdlttf.Font
	texture   *sdl.Texture
	surface   *sdl.Surface

	song *Music

	frame   = time.NewTicker(1000 / Fps * time.Millisecond)
	profile *os.File
)

func main() {
	runtime.LockOSThread()
	log.SetFlags(0)
	rand.Seed(time.Now().UnixNano())
	config.Parse()
	Profile()
	InitSDL()
	defer Quit()
	Load()
	Loop()
}

func InitSDL() {
	log.SetPrefix("sdl: ")

	err := sdl.Init(sdl.INIT_EVERYTHING &^ sdl.INIT_AUDIO)
	if err != nil {
		log.Fatal(err)
	}

	err = sdlttf.Init()
	if err != nil {
		log.Fatal(err)
	}

	err = sdl.InitSubSystem(sdl.INIT_AUDIO)
	if err != nil {
		log.Print(err)
	}

	err = sdlmixer.OpenAudio(44100, sdl.AUDIO_S16, 2, 4096)
	if err != nil {
		log.Print(err)
	}

	mxflag := sdlmixer.INIT_OGG
	nxflag, err := sdlmixer.Init(mxflag)
	if err != nil {
		log.Print(err)
	} else if nxflag != mxflag {
		log.Print("failed to initialize support for OGG")
	}
	sdlmixer.AllocateChannels(128)

	wflag := sdl.WINDOW_RESIZABLE
	if config.Fullscreen {
		wflag |= sdl.WINDOW_FULLSCREEN_DESKTOP
	}

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "linear")

	screen, err = newDisplay(W, H, wflag)
	if err != nil {
		log.Fatal(err)
	}
	screen.SetLogicalSize(W, H)

	sdl.ShowCursor(0)

	texture, err = screen.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, W, H)
	if err != nil {
		log.Fatal(err)
	}

	surface, err = sdl.CreateRGBSurfaceWithFormat(sdl.SWSURFACE, W, H, 32, sdl.PIXELFORMAT_ABGR8888)
	if err != nil {
		log.Fatal(err)
	}

	screen.SetTitle("Trip on the Funny Boat")
	screen.SetDrawColor(sdlcolor.Black)
	screen.Clear()
	screen.Present()

	icon, err := LoadSurface("kuvake")
	if err != nil {
		log.Print(err)
	} else {
		screen.SetIcon(icon)
	}
	icon.Free()

	sdl.StartTextInput()
}

func Load() {
	smallFont = LoadFont("Vera", 14)
	bigFont = LoadFont("Vera", 24)

	song = LoadMusic("JDruid-Trip_on_the_Funny_Boat")
	song.Play()

	InitWater()
	InitClouds()
}

func Loop() {
	var menu Menu
	var options Options
	var highscores Highscores
	var game Game

	mainMode := []string{"New Game", "High Scores", "Options", "Quit"}
	playMode := []string{"Story Mode", "Endless Mode"}

	menu.Init()
	options.Init()
	highscores.Init()
	game.Init()

	mainSelection := 0
loop:
	for {
		mainSelection = menu.Run(mainMode, mainSelection)
		switch mainSelection {
		case 0, 1: // Play, High Scores
			selection := 0
			selection = menu.Run(playMode, selection)

			var endless bool
			switch selection {
			case 0:
				endless = false
			case 1:
				endless = true
			default:
				continue loop
			}

			score := -1
			if mainSelection == 0 {
				score = game.Run(endless)
			}
			highscores.Run(endless, score)
		case 2: // Options
			options.Run()
		default: // Quit
			break loop
		}
	}
}

func Quit() {
	sdl.Quit()
	if profile != nil {
		pprof.StopCPUProfile()
		profile.Close()
	}
}

func Profile() {
	if config.Profile == "" {
		return
	}

	var err error

	log.SetPrefix("profile: ")
	profile, err = os.Create(config.Profile)
	if err != nil {
		log.Println(err)
		return
	}
	pprof.StartCPUProfile(profile)
}
