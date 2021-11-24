package ansi

import "fmt"

type Log struct {
	Message string
	Color   Color
	Prefix  string
	Suffix  string
}

func (l Log) String() string {
	return l.Prefix + l.Color.Fg() + l.Message + Reset + l.Suffix + "\n"
}

func NewLog(message string, a ...interface{}) *Log {
	return &Log{
		Message: fmt.Sprintf(message, a...),
		Color:   White,
		Prefix:  Blue.Fg() + "  ==> ",
		Suffix:  Reset,
	}
}

func (l *Log) SetColor(c Color) *Log {
	l.Color = c
	return l
}

func (l *Log) SetPrefix(p string) *Log {
	l.Prefix = p
	return l
}

func (l Log) Render() {
	fmt.Print(l.String())
}
