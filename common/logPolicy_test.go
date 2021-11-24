package common_test

import (
	"log/syslog"
	"testing"

	"github.com/wormable/nest/common"
)

func TestMatch(t *testing.T) {
	defaultPolicy := common.LogRule{}

	if !(defaultPolicy.Match(syslog.LOG_DEBUG) && defaultPolicy.Match(syslog.LOG_CRIT) && defaultPolicy.Match(syslog.LOG_ERR)) {
		t.Errorf("an empty rule should always log everything")
	}

	simplePolicy := common.LogRule{
		Level: "error",
	}

	if !(simplePolicy.Match(syslog.LOG_DEBUG) == false && simplePolicy.Match(syslog.LOG_ERR) == true && simplePolicy.Match(syslog.LOG_CRIT) == true) {
		t.Errorf("error in simple rule (weird)")
	}

	matrix := []struct {
		code   string
		input  syslog.Priority
		output bool
	}{
		{"level > error", syslog.LOG_DEBUG, false},
		{"level > error", syslog.LOG_ERR, false},
		{"level > error", syslog.LOG_CRIT, true},

		{"error > level", syslog.LOG_DEBUG, true},
		{"error > level", syslog.LOG_ERR, false},
		{"error > level", syslog.LOG_CRIT, false},

		{"level == info", syslog.LOG_INFO, true},
		{"level == info", syslog.LOG_DEBUG, false},

		{"info == level", syslog.LOG_INFO, true},
		{"info == level", syslog.LOG_DEBUG, false},

		{"level != info", syslog.LOG_INFO, false},
		{"level != info", syslog.LOG_DEBUG, true},

		{"info != level", syslog.LOG_INFO, false},
		{"info != level", syslog.LOG_DEBUG, true},

		{"level < warning", syslog.LOG_DEBUG, true},
		{"level < warning", syslog.LOG_WARNING, false},
		{"level < warning", syslog.LOG_CRIT, false},

		{"warning < level", syslog.LOG_DEBUG, false},
		{"warning < level", syslog.LOG_WARNING, false},
		{"warning < level", syslog.LOG_CRIT, true},

		{"info >= level", syslog.LOG_DEBUG, true},
		{"info >= level", syslog.LOG_INFO, true},
		{"info >= level", syslog.LOG_CRIT, false},

		{"level >= info", syslog.LOG_DEBUG, false},
		{"level >= info", syslog.LOG_INFO, true},
		{"level >= info", syslog.LOG_CRIT, true},

		{"info <= level", syslog.LOG_DEBUG, false},
		{"info <= level", syslog.LOG_INFO, true},
		{"info <= level", syslog.LOG_CRIT, true},

		{"level <= info", syslog.LOG_DEBUG, true},
		{"level <= info", syslog.LOG_INFO, true},
		{"level <= info", syslog.LOG_CRIT, false},
	}

	var rule common.LogRule
	for _, child := range matrix {
		rule.When = child.code

		if rule.Match(child.input) == !child.output {
			t.Errorf("%s should return %t with ( %s ), returned %t", child.code, child.output, common.Levels[child.input], !child.output)
		}
	}
}
