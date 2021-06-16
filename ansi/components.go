package ansi

import (
	"fmt"
)

func Text(message string, foreground RGB) {
	Println(message, []Modifier{
		Color{
			Kind:   Foreground,
			Values: foreground,
		},
	})
}

func Overline(message string, background RGB) {
	Println(" "+message+" ", []Modifier{
		Color{Kind: Background, Values: background},
		Color{Kind: Foreground, Values: Black},
		Bold,
	})
}

func Space() {
	fmt.Println()
}

func Spaces(n int) {
	for i := 0; i < n; i++ {
		Space()
	}
}

func Block(message string, background RGB) {
	Println(Padding(PaddingSizes{
		Content: message,
		Top:     1,
		Bottom:  1,
		Right:   8,
		Left:    8,
	}), []Modifier{
		Color{Kind: Background, Values: background},
		Color{Kind: Foreground, Values: White},
	})
}