package ansi

import "os"

var PrintAnsi = true

var Reset string
var Bold string
var Dim string
var Underline string
var Italic string
var Blink string
var Reverse string
var Hidden string
var StrikeThrough string

func init() {
	for _, arg := range os.Args {
		if arg == "--no-ansi" {
			PrintAnsi = false
			break
		}
	}

	if PrintAnsi {
		Reset = "\033[0m"
		Bold = "\033[1m"
		Dim = "\033[2m"
		Italic = "\033[3m"
		Underline = "\033[4m"
		Blink = "\033[5m"
		Reverse = "\033[7m"
		Hidden = "\033[8m"
		StrikeThrough = "\033[9m"
	}
}
