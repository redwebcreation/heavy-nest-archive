package common

import (
	"fmt"
	"strings"
)

type LogPolicy struct {
	When  string
	Level string
}

const (
	DebugLevel = iota
	InfoLevel
	WarningLevel
	ErrorLevel
	FatalLevel
)

var Levels = []string{"debug", "info", "warning", "error", "fatal"}

// The order matters here
var operators = []string{"<=", ">=", "<", ">", "==", "!="}

func getLevelValue(level string) int {
	for k, cmp := range Levels {
		if cmp == level {
			return k
		}
	}

	return -1
}

func (l LogPolicy) MustCompile(level int) (bool, error) {
	if 0 > level || level > len(Levels) {
		return false, fmt.Errorf("log level must be between 0 and %d, given %d", len(Levels), level)
	}

	if l.When == "" {
		if l.Level == "" {
			return true, nil
		}

		return level >= getLevelValue(l.Level), nil
	}

	expr := strings.ReplaceAll(l.When, "level", Levels[level])

	for _, condition := range strings.Split(expr, "||") {
		for _, op := range operators {
			if strings.Contains(condition, op) {
				terms := strings.Split(condition, op)
				leftTerm := getLevelValue(strings.TrimSpace(terms[0]))
				rightTerm := getLevelValue(strings.TrimSpace(terms[len(terms)-1]))

				result := map[string]bool{
					">":  leftTerm > rightTerm,
					">=": leftTerm >= rightTerm,
					"<":  leftTerm < rightTerm,
					"<=": leftTerm <= rightTerm,
					"==": leftTerm == rightTerm,
					"!=": leftTerm != rightTerm,
				}[op]

				if result {
					return true, nil
				}

				break
			}
		}
	}

	return false, nil
}

func (l LogPolicy) ShouldLog(level int) bool {
	log, _ := l.MustCompile(level)
	return log
}

func (l LogPolicy) Log(level int, message string, context ...string) {
	if !l.ShouldLog(level) {
		return
	}
}
