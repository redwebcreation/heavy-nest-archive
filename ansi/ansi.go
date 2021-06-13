package ansi

import (
	"fmt"
	"strings"
)

var Black = RGB{
	Red:   0,
	Green: 0,
	Blue:  0,
}

var White = RGB{
	Red:   255,
	Green: 255,
	Blue:  255,
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
var Blue = RGB{
	Red:   59,
	Green: 130,
	Blue:  246,
}
var Purple = RGB{
	Red:   139,
	Green: 92,
	Blue:  246,
}

func Print(message string, modifiers []Modifier) {
	lines := strings.Split(message, "\n")

	for _, line := range lines {
		for _, modifier := range modifiers {
			fmt.Print(modifier.String())
		}

		fmt.Print(line + "\033[0m")
	}
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
