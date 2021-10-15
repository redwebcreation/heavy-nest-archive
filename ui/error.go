package ui

import (
	"fmt"
	"os"
)

func Check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, Red.AsFg()+"Error: "+err.Error()+Stop)
		os.Exit(1)
	}
}
