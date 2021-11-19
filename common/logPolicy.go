package common

import (
	"encoding/json"
	"fmt"
	"github.com/wormable/ui"
	"log"
	"log/syslog"
	"os"
	"strings"
	"time"
)

type Fields map[string]string

type LogRedirection struct {
	Type     string `json:"type"`
	Path     string `json:"path,omitempty"`
	Facility string `json:"facility,omitempty"`
	Format   string `json:"format"`
}

type LogRule struct {
	When         string           `json:"when"`
	Level        string           `json:"level"`
	Redirections []LogRedirection `json:"redirections"`
	Format       string           `json:"format"`
}

type LogPolicy struct {
	Name   string    `json:"name"`
	Rules  []LogRule `json:"rules"`
	Format string    `json:"format"`

	context Fields
}

var Levels = []string{"emerg", "alert", "crit", "err", "warning", "notice", "info", "debug"}

// The order matters here <= must be before < and >= must be before >
var operators = []string{"<=", ">=", "<", ">", "==", "!="}

func getLevelValue(level string) syslog.Priority {
	for k, cmp := range Levels {
		if cmp == level {
			return syslog.Priority(k)
		}
	}

	return -1
}

func (l LogRule) ShouldLog(level syslog.Priority) (bool, error) {
	if 0 > level || int(level) > len(Levels) {
		return false, fmt.Errorf("log level must be between 0 and %d, given %d", len(Levels), level)
	}

	if l.When == "" {
		if l.Level == "" {
			return true, nil
		}

		return getLevelValue(l.Level) >= level, nil
	}
	expr := strings.ReplaceAll(l.When, "level", Levels[level])

	for _, condition := range strings.Split(expr, "||") {
		for _, op := range operators {
			if strings.Contains(condition, op) {
				terms := strings.Split(condition, op)
				lhs := getLevelValue(strings.TrimSpace(terms[0]))
				rhs := getLevelValue(strings.TrimSpace(terms[len(terms)-1]))

				result := map[string]bool{
					">":  lhs > rhs,
					">=": lhs >= rhs,
					"<":  lhs < rhs,
					"<=": lhs <= rhs,
					"==": lhs == rhs,
					"!=": lhs != rhs,
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

func (l LogRule) Match(level syslog.Priority) bool {
	compiled, _ := l.ShouldLog(level)
	return compiled
}

func (l LogRule) Log(format string, level syslog.Priority, message string, context Fields) {
	if format == "" {
		format = l.Format
	}

	context["level"] = Levels[level]
	context["time"] = time.Now().Format("2006/01/02 15:04:05")
	context["message"] = message

	for _, redirection := range l.Redirections {
		if redirection.Type == "stdout" {
			log.SetOutput(os.Stdout)
		} else if redirection.Type == "stderr" {
			log.SetOutput(os.Stderr)
		} else if redirection.Type == "file" {
			logFile, err := os.OpenFile(redirection.Path, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
			ui.Check(err)

			defer logFile.Close()

			log.SetOutput(logFile)
		} else if redirection.Type == "syslog" {
			writer, err := syslog.New(getLevelValue(l.Level), "")
			ui.Check(err)

			log.SetOutput(writer)
		} else {
			log.SetOutput(os.Stdout)
		}

		out, _ := json.Marshal(context)
		log.Printf("%s", out)
	}
}

func (l LogPolicy) WithContext(context Fields) LogPolicy {
	l.context = context

	return l
}

func (l LogPolicy) Log(level syslog.Priority, message string) {
	for _, rule := range l.Rules {
		if rule.Match(level) {
			if l.context == nil {
				l.context = make(Fields, 3)
			}

			rule.Log(l.Format, level, message, l.context)
		}
	}
}
