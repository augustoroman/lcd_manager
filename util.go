package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"code.google.com/p/sadbox/color"
)

type RGB struct{ R, G, B uint8 }

func (r RGB) MarshalText() ([]byte, error) {
	val := fmt.Sprintf(`%s`, color.RGBToHex(r.R, r.G, r.B))
	return []byte(val), nil
}
func (r *RGB) UnmarshalText(p []byte) error {
	r.R, r.G, r.B = color.HexToRGB(color.Hex(p))
	return nil
}

type ByteString []byte

func (l ByteString) MarshalText() ([]byte, error) { return json.Marshal(string(l)) }
func (l *ByteString) UnmarshalText(p []byte) error {
	var s string
	err := json.Unmarshal(p, &s)
	*l = []byte(s)
	return err
}

// dropN drops the number of bytes for functions return that and just returns
// the error.  For example:
//
//   err := dropN(writer.Write(data))
//
func dropN(n int, e error) error { return e }

func asByte(val string) uint8 { n, _ := strconv.ParseUint(val, 10, 8); return uint8(n) }
func asBool(val string) bool  { return val == "true" }
func asColor(val string) RGB {
	r, g, b := color.HexToRGB(color.Hex(val))
	return RGB{r, g, b}
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

func unquote(line string) string {
	if res, err := strconv.Unquote(`"` + line + `"`); err == nil {
		// log.Printf("Unquoted: [%s] -> %q", line, res)
		return res
		// } else {
		//  log.Printf("Error unquoting: [%s]: %v", line, err)
	}
	return line
}
