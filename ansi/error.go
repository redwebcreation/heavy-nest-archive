package ansi

import (
	"fmt"
	"os"
)

func Check(err error) {
	if err == nil {
		return
	}

	if !PrintAnsi {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	} else {
		_, _ = fmt.Fprintf(os.Stderr, Red.Fg()+"Error: %s\n"+Reset, err)
	}

	os.Exit(1)
}

func IfErr(err error, handler func(err string, printAnsi bool) string) {
	if err == nil {
		return
	}

	_, _ = fmt.Fprintf(os.Stderr, handler(err.Error(), PrintAnsi))
	os.Exit(1)
}

func IfErrPanic(err error, handler func(err string, printAnsi bool) string) {
	if err == nil {
		return
	}

	panic(handler(err.Error(), PrintAnsi))
}
