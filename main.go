package main

import (
	"flag"
	"image/color"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/augustoroman/multierror"
	"github.com/pwaller/go-hexcolor"

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

func main() {
	port := flag.String("port", "/dev/tty.usbmodem1451", "COM port that LCD is on.")
	baud := flag.Int("baud", 9600, "Baud rate to communicate at.")
	addr := flag.String("addr", ":12000", "Web address to bind to.")
	flag.Parse()

	var lcd LCD
	if *port == "" {
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
			bgcolor:    color.RGBA{200, 200, 50, 0},
			contrast:   200,
			brightness: 180,
			on:         true,
		},
		lcd: lcd,
	}
	s.configure(16, 2)
	lcd.SetSize(16, 2)

	m := martini.Classic()
	m.Handlers(martini.Recovery())
	m.Get("/", http.FileServer(rice.MustFindBox("www").HTTPBox()).ServeHTTP)
	m.Post("/set", s.Set)
	http.ListenAndServe(*addr, m)
}

type settings struct {
	display    [][]byte
	bgcolor    color.RGBA
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

type server struct {
	settings

	lines   []string
	linePos []int

	lcd LCD
}

func asByte(val string) uint8 { n, _ := strconv.ParseUint(val, 10, 8); return uint8(n) }
func asBool(val string) bool  { return val == "true" }
func asColor(val string) color.RGBA {
	r, g, b, a := hexcolor.HexToRGBA(hexcolor.Hex(val))
	return color.RGBA{r, g, b, a}
}

func (s *server) Update() {
	s.render()
	s.settings.apply(s.lcd)
}
func (s *server) SetLines(lines []string) {
	s.lines = lines
	s.linePos = make([]int, len(lines))
	s.Update()
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
		writeline(slice(s.lines[i]+buffer, s.linePos[i]), s.display[i])
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
		case "background":
			s.bgcolor = asColor(val)
		case "on":
			s.on = asBool(val)
		case "lines":
			s.SetLines(vals)
		}
	}
}
