package main

import (
	"flag"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"code.google.com/p/sadbox/color"
	"github.com/augustoroman/multierror"

	"github.com/GeertJohan/go.rice"
	"github.com/augustoroman/serial_lcd"
	"github.com/go-martini/martini"
)

type LCD interface {
	SetBG(r, g, b uint8) error
	SetOn(on bool) error
	SetBrightness(b uint8) error
	SetContrast(c uint8) error
	SetAutoscroll(on bool) error
	SetSize(cols, rows uint8) error
	Clear() error
	Home() error
	MoveTo(col, row uint8) error
	MoveForward() error
	MoveBack() error
	io.Writer
}

type FakeLcd struct{}

func (f FakeLcd) SetBG(r, g, b uint8) error      { return nil }
func (f FakeLcd) SetOn(on bool) error            { return nil }
func (f FakeLcd) SetBrightness(b uint8) error    { return nil }
func (f FakeLcd) SetContrast(c uint8) error      { return nil }
func (f FakeLcd) SetAutoscroll(on bool) error    { return nil }
func (f FakeLcd) SetSize(cols, rows uint8) error { return nil }
func (f FakeLcd) Clear() error                   { return nil }
func (f FakeLcd) Home() error                    { return nil }
func (f FakeLcd) MoveTo(col, row uint8) error    { return nil }
func (f FakeLcd) MoveForward() error             { return nil }
func (f FakeLcd) MoveBack() error                { return nil }
func (f FakeLcd) Write(b []byte) (int, error)    { return len(b), nil }

type RGB struct{ R, G, B uint8 }

func main() {
	port := flag.String("port", "/dev/serial/by-id/usb-239a_Adafruit_Industries-if00", "COM port that LCD is on.")
	baud := flag.Int("baud", 9600, "Baud rate to communicate at.")
	addr := flag.String("addr", ":12000", "Web address to bind to.")
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
		settings: settings{
			bgcolor:    RGB{200, 200, 50},
			contrast:   200,
			brightness: 180,
			on:         true,
		},
		lcd: lcd,
		ch:  make(chan server),
	}
	go lcdLoop(s.ch)
	s.configure(16, 2)
	lcd.SetSize(16, 2)
	s.SetLines("", "")
	s.Update()

	m := martini.Classic()
	m.Handlers(martini.Recovery())
	m.Get("/", http.FileServer(rice.MustFindBox("www").HTTPBox()).ServeHTTP)
	m.Post("/set", s.Set)
	http.ListenAndServe(*addr, m)
}

type settings struct {
	display    [][]byte
	bgcolor    RGB
	contrast   uint8
	brightness uint8
	on         bool
}

func (s settings) apply(lcd LCD) error {
	var errs multierror.Accumulator
	errs.Push(lcd.SetOn(s.on))
	errs.Push(lcd.SetBG(s.bgcolor.R, s.bgcolor.G, s.bgcolor.B))
	errs.Push(lcd.SetBrightness(s.brightness))
	errs.Push(lcd.SetContrast(s.contrast))
	errs.Push(lcd.SetOn(s.on))
	errs.Push(lcd.Home())
	errs.Push(dropN(lcd.Write(s.display[0])))
	errs.Push(dropN(lcd.Write(s.display[1])))
	return errs.Error()
}

func (s *settings) configure(width, height int) {
	s.display = make([][]byte, height)
	for i := 0; i < height; i++ {
		s.display[i] = make([]byte, width)
	}
}

func dropN(n int, e error) error { return e }

func lcdLoop(newSettings chan server) {
	s := <-newSettings
	timer := time.NewTicker(100 * time.Millisecond)
	defer timer.Stop()
	last := time.Now()
	var open bool
	for {
		select {
		case s, open = <-newSettings:
			if !open {
				return
			}
		case t := <-timer.C:
			s.advance(t.Sub(last))
		}
		s.apply(s.lcd)
		last = time.Now()
	}
}

func (s *server) advance(dt time.Duration) {
	if s.rainbow {
		h, ss, v := color.RGBToHSV(s.bgcolor.R, s.bgcolor.G, s.bgcolor.B)
		h = math.Mod(h+dt.Seconds()/10.0, 1.0)
		s.bgcolor.R, s.bgcolor.G, s.bgcolor.B = color.HSVToRGB(h, ss, v)
	}
}

type server struct {
	settings

	lines   []string
	linePos []float64

	lcd LCD

	rainbow bool

	ch chan server
}

func asByte(val string) uint8 { n, _ := strconv.ParseUint(val, 10, 8); return uint8(n) }
func asBool(val string) bool  { return val == "true" }
func asColor(val string) RGB {
	r, g, b := color.HexToRGB(color.Hex(val))
	return RGB{r, g, b}
}

func (s *server) Update() error {
	s.render()
	s.ch <- *s
	return nil
	//return s.settings.apply(s.lcd)
}
func (s *server) SetLines(lines ...string) {
	s.lines = lines
	s.linePos = make([]float64, len(lines))
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func slice(line string, start int) string {
	if len(line) < start {
		return ""
	}
	return line[start:]
}
func writeline(line string, dest []byte) {
	for i := 0; i < min(len(line), len(dest)); i++ {
		dest[i] = line[i]
	}
	for i := len(line); i < len(dest); i++ {
		dest[i] = ' '
	}
}
func (s *server) render() {
	const buffer = "   "
	for i := 0; i < min(len(s.lines), len(s.display)); i++ {
		writeline(slice(s.lines[i]+buffer, int(s.linePos[i])), s.display[i])
	}
	for i := len(s.lines); i < len(s.display); i++ {
		writeline("", s.display[i])
	}
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
		case "brightness":
			s.brightness = asByte(val)
		case "contrast":
			s.contrast = asByte(val)
		case "color":
			s.bgcolor = asColor(val)
		case "on":
			s.on = asBool(val)
		case "rainbow":
			s.rainbow = asBool(val)
		case "line[]":
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
