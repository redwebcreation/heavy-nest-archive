package ansi

import "strings"

type Color [3]uint8

var Blue Color = [3]uint8{59, 130, 246}
var White Color = [3]uint8{255, 255, 255}
var Gray Color = [3]uint8{147, 148, 153}
var Red Color = [3]uint8{239, 68, 68}
var Yellow Color = [3]uint8{245, 158, 11}
var Teal Color = [3]uint8{56, 178, 172}
var Green Color = [3]uint8{16, 185, 129}

func FromHex(hex string) Color {
	var c Color
	hex = strings.TrimLeft(hex, "#")

	for i := 0; i < 3; i++ {
		c[i] = hex[i*2]*16 + hex[i*2+1]
	}
	return c
}

func New(r, g, b uint8) Color {
	return Color{r, g, b}
}

func (c Color) Fg() string {
	if !PrintAnsi {
		return ""
	}

	return "\033[38;2;" + string(c[0]) + ";" + string(c[1]) + ";" + string(c[2]) + "m"
}

func (c Color) Bg() string {
	if !PrintAnsi {
		return ""
	}

	return "\033[48;2;" + string(c[0]) + ";" + string(c[1]) + ";" + string(c[2]) + "m"
}
