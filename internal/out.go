package internal

import (
	"fmt"
	"os"
)

type Color [3]uint8
type Log struct {
	Message     string
	arrowString string
	arrow       *Color
	color       *Color
	nesting     int
	top         int
}

var Primary Color = [3]uint8{59, 130, 246}
var White Color = [3]uint8{255, 255, 255}
var Gray Color = [3]uint8{147, 148, 153}
var Red Color = [3]uint8{239, 68, 68}
var Green Color = [3]uint8{16, 185, 129}
var Stop = "\033[0m"
var Bold = "\033[1m"

func (c Color) AsFg() string {
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", c[0], c[1], c[2])
}

func (l Log) String() string {
	if l.arrow == nil {
		l.arrow = &Primary
	}

	if l.color == nil {
		l.color = &White
	}

	if l.arrowString == "" {
		l.arrowString = "==>"
	}

	str := ""

	for i := 0; i < int(l.nesting+1); i++ {
		str += "    " // 4 spaces, 1 tab
	}

	str += fmt.Sprintf("%s%s%s", l.arrow.AsFg()+Bold, l.arrowString, Stop)
	str += fmt.Sprintf(" %s%s%s", l.color.AsFg()+Bold, l.Message, Stop)

	if l.top != 0 {
		str = fmt.Sprintf("\033[%dA\033[K", l.top+1) + str
	}
	return str + "\n"
}

func (l Log) Print() {
	fmt.Println(l.String())
}

func Check(err error) {
	if err != nil {
		fmt.Println("\033[38;2;255;0;0" + err.Error() + "\033[0m")
		os.Exit(1)
	}
}

func NewLog(format string, a ...interface{}) Log {
	return Log{
		Message: fmt.Sprintf(format, a...),
	}
}

func (l Log) Arrow(c Color) Log {
	l.arrow = &c
	return l
}

func (l Log) Color(c Color) Log {
	l.color = &c
	return l
}

func (l Log) Top(n int) Log {
	l.top = n
	return l
}

func (l Log) Nesting(n int) Log {
	l.nesting = n
	return l
}

func (l Log) ArrowString(shape string) Log {
	l.arrowString = shape
	return l
}

func Title(format string, a ...interface{}) {
	fmt.Println(Bold + White.AsFg() + fmt.Sprintf(format, a...) + Stop)
	fmt.Println()
}
