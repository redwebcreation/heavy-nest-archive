package ansi

import (
	"fmt"
)

func CursorUp(n int) {
	if !PrintAnsi {
		return
	}

	if n == 1 {
		fmt.Printf("\033[A")
	} else {
		fmt.Printf("\033[%dA", n)
	}
}

func CursorDown(n int) {
	if !PrintAnsi {
		return
	}

	if n == 1 {
		fmt.Printf("\033[B")
	} else {
		fmt.Printf("\033[%dB", n)
	}
}

func CursorForward(n int) {
	if !PrintAnsi {
		return
	}

	if n == 1 {
		fmt.Printf("\033[C")
	} else {
		fmt.Printf("\033[%dC", n)
	}
}
func CursorBackward(n int) {
	if !PrintAnsi {
		return
	}

	if n == 1 {
		fmt.Printf("\033[D")
	} else {
		fmt.Printf("\033[%dD", n)
	}
}

func CursorNextLine(n int) {
	if !PrintAnsi {
		return
	}

	if n == 1 {
		fmt.Printf("\033[E")
	} else {
		fmt.Printf("\033[%dE", n)
	}
}

func CursorPreviousLine(n int) {
	if !PrintAnsi {
		return
	}

	if n == 1 {
		fmt.Printf("\033[F")
	} else {
		fmt.Printf("\033[%dF", n)
	}
}

func MoveCursor(x, y int) {
	if !PrintAnsi {
		return
	}

	fmt.Printf("\033[%d;%dH", y, x)
}

func MoveCursorHome() {
	if !PrintAnsi {
		return
	}

	fmt.Printf("\033[H")
}

func SaveCursorPosition() {
	if !PrintAnsi {
		return
	}

	fmt.Printf("\033 7")
}

func RestoreCursorPosition() {
	if !PrintAnsi {
		return
	}

	fmt.Printf("\033 8")
}


