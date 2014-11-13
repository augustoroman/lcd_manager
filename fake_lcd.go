package main

import "github.com/augustoroman/serial_lcd"

type FakeLcd struct{}

func (f FakeLcd) SetBG(r, g, b uint8) error      { return nil }
func (f FakeLcd) SetOn(On bool) error            { return nil }
func (f FakeLcd) SetBrightness(b uint8) error    { return nil }
func (f FakeLcd) SetContrast(c uint8) error      { return nil }
func (f FakeLcd) SetAutoscroll(On bool) error    { return nil }
func (f FakeLcd) SetSize(cols, rows uint8) error { return nil }
func (f FakeLcd) Clear() error                   { return nil }
func (f FakeLcd) Home() error                    { return nil }
func (f FakeLcd) MoveTo(col, row uint8) error    { return nil }
func (f FakeLcd) MoveForward() error             { return nil }
func (f FakeLcd) MoveBack() error                { return nil }
func (f FakeLcd) Write(b []byte) (int, error)    { return len(b), nil }

func (f FakeLcd) CreateCustomChar(spot uint8, c serial_lcd.Char) error { return nil }
