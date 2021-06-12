package core

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type ansiOutput bool

// Ansi Whether the terminal should print out Ansi characters.
var Ansi ansiOutput

func init() {
	for i, arg := range os.Args {
		if arg == "--no-ansi" {
			// Removes the argument so Cobra doesn't freak out.
			os.Args = os.Args[:i+copy(os.Args[i:], os.Args[i+1:])]
			Ansi = false
			break
		}
	}

	Ansi = true
}

func (ansi ansiOutput) Print(elements ...interface{}) {
	re := regexp.MustCompile("\033\\[[0-9;]+m(.+)\033\\[0m")

	for _, element := range elements {
		switch stringElement := element.(type) {
		case string:
			if !ansi {
				stringElement = re.ReplaceAllString(stringElement, "$1")
			}

			fmt.Print(stringElement)
		default:
			fmt.Print(element)
		}

	}
}

func (ansi ansiOutput) Println(elements ...interface{}) {
	for _, element := range elements {
		ansi.Print(element)
		ansi.NewLine()
	}
}

func (ansi ansiOutput) Printf(format string, elements ...interface{}) {
	re := regexp.MustCompile("\033\\[[0-9;]+m(.+)\033\\[0m")

	if !ansi {
		format = re.ReplaceAllString(format, "$1")
	}

	fmt.Printf(format, elements...)
}

func (ansi ansiOutput) Reset() {
	fmt.Print("\033[0m")
}

func (ansi ansiOutput) Success(message string) {
	ansi.Println("\033[32m" + message)
	ansi.Reset()
}

func (ansi ansiOutput) Warning(message string) {
	ansi.Println("\033[33m" + message)
	ansi.Reset()
}

func (ansi ansiOutput) Error(message string) {
	ansi.Println("\033[31m" + message)
	ansi.Reset()
}

func (ansi ansiOutput) Check(err error) {
	if err != nil {
		ansi.Error(err.Error())
	}
}

func (ansi ansiOutput) Block(message string, style string) {
	paddingX := 8
	paddingY := 1

	if !ansi {
		if !ansi {
			paddingX = 0
			paddingY = 0
		}

		line := style + strings.Repeat(" ", len(message)+paddingX*2) + "\033[0m\n"

		message = style + strings.Repeat(" ", paddingX) + message + strings.Repeat(" ", paddingX) + "\033[0m\n"

		topLines := strings.Repeat(line, paddingY)
		bottomLines := strings.Repeat(line, paddingY)

		ansi.Print(topLines + message + bottomLines)
	}
}

func (ansi ansiOutput) ErrorBlock(message string) {
	ansi.Block(message, "\033[1;97;41m")
}

func (ansi ansiOutput) WarningBlock(message string) {
	ansi.Block(message, "\033[1;97;43m")
}

func (ansi ansiOutput) SuccessBlock(message string) {
	ansi.Block(message, "\033[1;97;42m")
}

func (ansi ansiOutput) ClearLine() {
	if ansi {
		fmt.Print("\033[1A\033[K")
	}
}

func (ansi ansiOutput) StatusLoader(messages string, count *int) {
	ansi.ClearLine()
	ansi.Print(messages + strings.Repeat(".", *count) + "\n")

	*count++

	if *count > 3 {
		*count = 0
	}

}

func (ansi ansiOutput) NewLine() {
	fmt.Print("\n")
}
