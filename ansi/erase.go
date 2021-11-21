package ansi

import "fmt"

func ClearFromCursorToEnd() {
	if !PrintAnsi {
		return
	}
	fmt.Print("\033[0J")
}

func ClearFromStartToCursor() {
	if !PrintAnsi {
		return
	}
	fmt.Print("\033[1J")
}

func Clear() {
	if !PrintAnsi {
		return
	}
	fmt.Print("\033[2J")
}


func ClearLineFromCursorToEnd() {
	if !PrintAnsi {
		return
	}
	fmt.Print("\033[0K")
}

func ClearLineFromStartToCursor() {
	if !PrintAnsi {
		return
	}
	fmt.Print("\033[1K")
}

func ClearLine() {
	if !PrintAnsi {
		return
	}
	fmt.Print("\033[2K")
}