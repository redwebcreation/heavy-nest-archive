package ui

import "fmt"

func Title(format string, a ...interface{}) {
	fmt.Println(Bold + Primary.AsFg() + fmt.Sprintf(format, a...) + Stop)
	fmt.Println()
}
