package ansi

import (
	"fmt"
	"strings"
)
const (
	Foreground = "38;2;%d;%d;%d"
)

type RGB struct {
	Red   int
	Green int
	Blue  int
}

type Modifier interface {
	String() string
}

type Color struct {
	Kind   string
	Values RGB
}

func (color Color) String() string {
	return "\033[" + fmt.Sprintf(color.Kind, color.Values.Red, color.Values.Green, color.Values.Blue) + "m"
}

type Content string

func (content Content) String() string {
	return string(content)
}

type Effect string

func (effect Effect) String() string {
	return "\033[" + string(effect)
}

var Red = RGB{
	Red:   220,
	Green: 38,
	Blue:  38,
}
var Orange = RGB{
	Red:   245,
	Green: 158,
	Blue:  11,
}
var Green = RGB{
	Red:   16,
	Green: 185,
	Blue:  129,
}

func Println(message string, modifiers []Modifier) {
	lines := strings.Split(message, "\n")

	for _, line := range lines {
		for _, modifier := range modifiers {
			fmt.Print(modifier.String())
		}

		fmt.Print(line + "\033[0m\n")
	}
}

func Text(message string, foreground RGB) {
	Println(message, []Modifier{
		Color{
			Kind:   Foreground,
			Values: foreground,
		},
	})
}