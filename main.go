package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"time"

	"code.google.com/p/sadbox/color"
	"github.com/GeertJohan/go.rice"
	"github.com/augustoroman/multierror"
	"github.com/augustoroman/serial_lcd"
	"github.com/cactus/gostrftime"
	"github.com/go-martini/martini"
)

type LCD interface {
	SetBG(r, g, b uint8) error
	SetOn(on bool) error
	SetBrightness(b uint8) error
	SetContrast(c uint8) error
	SetAutoscroll(On bool) error
	SetSize(cols, rows uint8) error
	Clear() error
	Home() error
	MoveTo(col, row uint8) error
	MoveForward() error
	MoveBack() error
	CreateCustomChar(spot uint8, c serial_lcd.Char) error
	io.Writer
}

func main() {
	port := flag.String("port", "/dev/serial/by-id/usb-239a_Adafruit_Industries-if00", "COM port that LCD is On.")
	baud := flag.Int("baud", 9600, "Baud rate to communicate at.")
	addr := flag.String("addr", ":12000", "Web address to bind to.")
	settingsFilename := flag.String("settings", ".lcd_manager.settings", "Settings file.")
	flag.Parse()

	var lcd LCD
	if *port == "" {
		log.Println("Using fake LCD interface since empty --port specified.")
		lcd = FakeLcd{}
	} else {
		var err error
		lcd, err = serial_lcd.Open(*port, *baud)
		if err != nil {
			log.Fatal(err)
		}
	}

	s := &server{
		Settings: Settings{
			BgColor:    RGB{200, 200, 50},
			Contrast:   200,
			Brightness: 180,
			On:         true,
		},
		settingsFile: *settingsFilename,
		lcd:          lcd,
		ch:           make(chan server),
		Frames:       map[string]string{},
	}
	if err := s.Load(); err != nil {
		log.Println(err)
	} else {
		log.Println("Loading settings from", *settingsFilename)
	}
	go lcdLoop(s.ch)
	s.configure(16, 2)
	lcd.SetSize(16, 2)
	s.SetLines(s.Lines...)

	// pacmanChars(s.lcd)
	// invaderChars(s.lcd)
	watcherChars(s.lcd)

	s.Update()

	m := martini.Classic()
	m.Handlers(martini.Recovery())
	m.Post("/set", s.Set)

	m.Get("/settings", s.GetSettings)
	m.Get("/frames", s.GetFrames)
	m.Put("/frames", s.SetFrames)

	m.Get("/**", http.FileServer(rice.MustFindBox("www").HTTPBox()).ServeHTTP)
	if err := http.ListenAndServe(*addr, m); err != nil {
		log.Fatal(err)
	}
}

type Settings struct {
	display    []ByteString
	BgColor    RGB
	Contrast   uint8
	Brightness uint8
	On         bool
}

func (s Settings) apply(lcd LCD) error {
	var errs multierror.Accumulator
	errs.Push(lcd.SetOn(s.On))
	errs.Push(lcd.SetBG(s.BgColor.R, s.BgColor.G, s.BgColor.B))
	errs.Push(lcd.SetBrightness(s.Brightness))
	errs.Push(lcd.SetContrast(s.Contrast))
	errs.Push(lcd.SetOn(s.On))
	errs.Push(lcd.Home())
	errs.Push(dropN(lcd.Write(s.display[0])))
	errs.Push(dropN(lcd.Write(s.display[1])))
	return errs.Error()
}

func (s *Settings) configure(width, height int) {
	s.display = make([]ByteString, height)
	for i := 0; i < height; i++ {
		s.display[i] = make(ByteString, width)
	}
}

func (s *server) Load() error {
	if data, err := ioutil.ReadFile(s.settingsFile); err == nil {
		if err := json.Unmarshal(data, &s); err != nil {
			return fmt.Errorf("Error loading settings: %v", err)
		}
	} else {
		return fmt.Errorf("Error reading settings file: %v", err)
	}
	return nil
}

func (s *server) Save() error {
	data, err := json.Marshal(s)
	if err == nil {
		return ioutil.WriteFile(s.settingsFile, data, 0644)
	}
	return err
}

func (s *server) GetSettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(s)
}

func lcdLoop(newSettings chan server) {
	s := <-newSettings
	timer := time.NewTicker(100 * time.Millisecond)
	defer timer.Stop()
	last := time.Now()
	var open bool
	for {
		select {
		case s, open = <-newSettings:
			s.Save()
			if !open {
				return
			}
		case t := <-timer.C:
			s.advance(t.Sub(last))
		}
		s.render()
		s.apply(s.lcd)
		last = time.Now()
	}
}

func (s *server) advance(dt time.Duration) {
	if s.Rainbow {
		h, ss, v := color.RGBToHSV(s.BgColor.R, s.BgColor.G, s.BgColor.B)
		h = math.Mod(h+dt.Seconds()/7.0, 1.0)
		s.BgColor.R, s.BgColor.G, s.BgColor.B = color.HSVToRGB(h, ss, v)
	}
}

type server struct {
	Settings

	Lines   []string
	LinePos []float64

	Frames map[string]string

	lcd LCD

	Rainbow bool

	ch           chan server
	settingsFile string
}

func (s *server) Update() error {
	s.ch <- *s
	return nil
	//return s.Settings.apply(s.lcd)
}
func (s *server) SetLines(lines ...string) {
	s.Lines = lines
	s.LinePos = make([]float64, len(lines))
}

func (s *server) render() {
	const buffer = "   "
	for i := 0; i < min(len(s.Lines), len(s.display)); i++ {
		line := slice(s.Lines[i]+buffer, int(s.LinePos[i]))
		line = unquote(line)
		line = gostrftime.Format(line, time.Now())
		writeline(line, s.display[i])
	}
	for i := len(s.Lines); i < len(s.display); i++ {
		writeline("", s.display[i])
	}
}

func (s *server) GetFrames(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(s.Frames)
}
func (s *server) SetFrames(w http.ResponseWriter, r *http.Request) {
	var frames map[string]string
	if err := json.NewDecoder(r.Body).Decode(&frames); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.Frames = frames
	s.Update()
}

func (s *server) Set(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for key, vals := range r.Form {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]
		switch key {
		case "Brightness":
			s.Brightness = asByte(val)
		case "Contrast":
			s.Contrast = asByte(val)
		case "BgColor":
			s.BgColor = asColor(val)
			// log.Println(s.BgColor)
		case "On":
			s.On = asBool(val)
		case "Rainbow":
			s.Rainbow = asBool(val)
		case "Lines[]":
			s.SetLines(vals...)
		default:
			log.Printf("Unknown form key %q = %q", key, vals)
		}
	}
	if err := s.Update(); err != nil {
		log.Printf("Failed to update lcd: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
