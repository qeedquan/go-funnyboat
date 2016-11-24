package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type Config struct {
	Dir           string
	Resource      string
	Name          string
	Profile       string
	Fullscreen    bool
	Invincibility bool
	Sound         bool
	Music         bool
	Particles     bool
}

func (c *Config) Parse() {
	var fullscreen, invincible, noSound, noMusic, noParticles bool

	userDir := "."
	user, err := user.Current()
	if err == nil {
		userDir = user.HomeDir
	}

	if runtime.GOOS != "windows" {
		userDir = filepath.Join(userDir, ".funnyboat")
	} else {
		userDir = filepath.Join(userDir, "Funny Boat")
	}

	flag.StringVar(&c.Dir, "c", userDir, "config directory")
	flag.StringVar(&c.Resource, "r", "data", "resource directory")
	flag.StringVar(&c.Profile, "p", "", "turn on profiling and output to file")
	flag.BoolVar(&fullscreen, "f", false, "fullscreen")
	flag.BoolVar(&invincible, "i", false, "invincible")
	flag.BoolVar(&noSound, "ns", false, "no sound")
	flag.BoolVar(&noMusic, "nm", false, "no music")
	flag.BoolVar(&noParticles, "np", false, "no particles")
	flag.Parse()

	c.Particles = true
	c.Sound = true
	c.Music = true
	c.Name = "Funny Boat"
	c.Load()

	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "f":
			c.Fullscreen = fullscreen
		case "i":
			c.Invincibility = invincible
		case "ns":
			c.Sound = !noSound
		case "nm":
			c.Music = !noMusic
		case "np":
			c.Particles = !noParticles
		}
	})
}

func (c *Config) Path() (string, error) {
	err := os.MkdirAll(c.Dir, 0755)
	if err != nil && !os.IsExist(err) {
		return "", err
	}

	return c.Dir, nil
}

func (c *Config) Filename() (string, error) {
	path, err := c.Path()
	if err != nil {
		return "", err
	}
	return filepath.Join(path, "config"), nil
}

func (c *Config) Load() {
	var err error

	defer func() {
		if err != nil {
			log.Print("load failure: ", err)
		}
	}()

	filename, err := c.Filename()
	if err != nil {
		return
	}

	log.SetPrefix("config: ")

	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	log.Print("load success: ", filename)

	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		tokens := strings.Split(line, "\t")
		if len(tokens) < 2 {
			continue
		}

		for i := range tokens {
			tokens[i] = strings.TrimSpace(tokens[i])
		}

		switch strings.ToLower(tokens[0]) {
		case "particles":
			c.Particles, _ = strconv.ParseBool(tokens[1])
		case "invincibility":
			c.Invincibility, _ = strconv.ParseBool(tokens[1])
		case "music":
			c.Music, _ = strconv.ParseBool(tokens[1])
		case "name":
			c.Name = tokens[1]
		case "sound":
			c.Sound, _ = strconv.ParseBool(tokens[1])
		}
	}
}

func (c *Config) Save() {
	var filename string
	var err error

	log.SetPrefix("config: ")
	defer func() {
		if err != nil {
			log.Print("save failure: ", err)
		} else {
			log.Printf("saved to %q", filename)
		}
	}()

	filename, err = c.Filename()
	if err != nil {
		return
	}

	f, err := os.Create(filename)
	if err != nil {
		return
	}

	w := bufio.NewWriter(f)
	fmt.Fprintf(w, "particles\t%v\n", c.Particles)
	fmt.Fprintf(w, "invincibility\t%v\n", c.Invincibility)
	fmt.Fprintf(w, "music\t%v\n", c.Music)
	fmt.Fprintf(w, "name\t%v\n", c.Name)
	fmt.Fprintf(w, "sound\t%v\n", c.Sound)

	flushErr := w.Flush()
	closeErr := f.Close()

	err = flushErr
	if err == nil {
		err = closeErr
	}
}
