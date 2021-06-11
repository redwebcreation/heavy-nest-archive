package ansi

import "fmt"

func Print(string string) {
	fmt.Println(string)
}

func Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func Error(message string) {
	Print("\033[31m" + message + "\033[0m")
}

func Check(err error) {
	if err != nil {
		Error(err.Error())
	}
}

func Success(message string) {
	Print("\033[32m" + message + "\033[0m")
}

func Warning(message string) {
	Print("\033[33m" + message + "\033[0m")
}
