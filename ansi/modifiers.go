package ansi

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	Foreground = "38;2;%d;%d;%d"
	Background = "48;2;%d;%d;%d"
)

type Modifier interface {
	String() string
}

type RGB struct {
	Red   int
	Green int
	Blue  int
}

type Color struct {
	Kind   string
	Values RGB
}

type PaddingSizes struct {
	Content string
	Top     int
	Bottom  int
	Right   int
	Left    int
}

type Content string

type Effect string

const Up = "A"
const Down = "B"
const Right = "C"
const Left = "D"

type Move struct {
	Direction string
	Repeat    int
}

var Reset Effect
var Bold Effect
var Dim Effect
var Italic Effect
var Underline Effect

func init() {
	Reset = "0"
	Bold = "1"
	Dim = "2"
	Italic = "3"
	Underline = "4"
}

func (effect Effect) String() string {
	return "\033[" + string(effect)
}

func (color Color) String() string {
	return "\033[" + fmt.Sprintf(color.Kind, color.Values.Red, color.Values.Green, color.Values.Blue) + "m"
}

func (content Content) String() string {
	return string(content)
}

func Padding(padding PaddingSizes) string {
	right := strings.Repeat(" ", padding.Right)
	left := strings.Repeat(" ", padding.Left)
	lines := strings.Split(padding.Content, "\n")

	for i, line := range lines {
		lines[i] = left + line + right
	}

	maxLineLength := arrayMax(lines)

	var paddedLines []string

	for i := 0; i < padding.Top; i++ {
		paddedLines = append(paddedLines, strings.Repeat(" ", maxLineLength))
	}

	for _, line := range lines {
		paddedLines = append(paddedLines, line)
	}

	for i := 0; i < padding.Bottom; i++ {
		paddedLines = append(paddedLines, strings.Repeat(" ", maxLineLength))
	}

	return strings.Join(paddedLines, "\n")
}

func (move Move) String() string {
	fmt.Println(Effect(strconv.Itoa(move.Repeat) + move.Direction).String())
	return Effect(strconv.Itoa(move.Repeat) + move.Direction).String()
}
