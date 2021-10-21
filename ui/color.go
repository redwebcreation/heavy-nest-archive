package ui

import (
	"fmt"
)

type Color [3]uint8

var Primary Color = [3]uint8{59, 130, 246}
var White Color = [3]uint8{255, 255, 255}
var Gray Color = [3]uint8{147, 148, 153}
var Red Color = [3]uint8{239, 68, 68}
var Yellow Color = [3]uint8{245,158,11}
var Teal Color = [3]uint8{56,178,172}
var Green Color = [3]uint8{16, 185, 129}
var Stop = "\033[0m"
var Bold = "\033[1m"
var Dim = "\033[2m"

func (c Color) Fg() string {
	return fmt.Sprintf("%s\033[38;2;%d;%d;%dm", Bold, c[0], c[1], c[2])
}

func (c Color) Bg() string {
	return fmt.Sprintf("\033[48;2;%d;%d;%dm", c[0], c[1], c[2])
}