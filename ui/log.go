package ui

import (
	"fmt"
)

type Log struct {
	Message     string
	arrowString string
	arrow       *Color
	color       *Color
	top         int
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

	str := "    "

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

func (l Log) ArrowString(shape string) Log {
	l.arrowString = shape
	return l
}
