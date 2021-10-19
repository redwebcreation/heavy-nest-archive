package ui

import "fmt"

func Title(format string, a ...interface{}) {
	fmt.Println(Primary.Fg() + fmt.Sprintf(format, a...) + Stop)
	fmt.Println()
}
