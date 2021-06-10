package ansi

import "fmt"

func Check(err error) {
	if err == nil {
		return
	}

	fmt.Println("\033[ " + err.Error() + "s")
}
