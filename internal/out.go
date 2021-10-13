package internal

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

var TermWidth int
var TermHeight int

func init() {
	w, h, err := term.GetSize(int(os.Stdin.Fd()))
	Check(err)
	TermWidth = w
	TermHeight = h
}

func Log(format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
}

func Check(err error) {
	if err != nil {
		fmt.Println("\033[38;2;255;0;0m" + err.Error() + "\033[0m")
		os.Exit(1)
	}
}
