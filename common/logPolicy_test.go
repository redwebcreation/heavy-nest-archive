package common_test

import (
	"testing"

	"github.com/wormable/nest/common"
)

func TestShouldLog(t *testing.T) {
	defaultPolicy := common.LogRule{}

	if !(defaultPolicy.ShouldLog(common.DebugLevel) && defaultPolicy.ShouldLog(common.FatalLevel) && defaultPolicy.ShouldLog(common.ErrorLevel)) {
		t.Errorf("an empty rule should always log everything")
	}

	simplePolicy := common.LogRule{
		Level: "error",
	}

	if !(simplePolicy.ShouldLog(common.DebugLevel) == false && simplePolicy.ShouldLog(common.ErrorLevel) == true && simplePolicy.ShouldLog(common.FatalLevel) == true) {
		t.Errorf("error in simple rule (weird)")
	}

	matrix := []struct {
		code   string
		input  int
		output bool
	}{
		{"level > error", common.DebugLevel, false},
		{"level > error", common.ErrorLevel, false},
		{"level > error", common.FatalLevel, true},

		{"error > level", common.DebugLevel, true},
		{"error > level", common.ErrorLevel, false},
		{"error > level", common.FatalLevel, false},

		{"level == info", common.InfoLevel, true},
		{"level == info", common.DebugLevel, false},

		{"info == level", common.InfoLevel, true},
		{"info == level", common.DebugLevel, false},

		{"level != info", common.InfoLevel, false},
		{"level != info", common.DebugLevel, true},

		{"info != level", common.InfoLevel, false},
		{"info != level", common.DebugLevel, true},

		{"level < warning", common.DebugLevel, true},
		{"level < warning", common.WarningLevel, false},
		{"level < warning", common.FatalLevel, false},

		{"warning < level", common.DebugLevel, false},
		{"warning < level", common.WarningLevel, false},
		{"warning < level", common.FatalLevel, true},

		{"info >= level", common.DebugLevel, true},
		{"info >= level", common.InfoLevel, true},
		{"info >= level", common.FatalLevel, false},

		{"level >= info", common.DebugLevel, false},
		{"level >= info", common.InfoLevel, true},
		{"level >= info", common.FatalLevel, true},

		{"info <= level", common.DebugLevel, false},
		{"info <= level", common.InfoLevel, true},
		{"info <= level", common.FatalLevel, true},

		{"level <= info", common.DebugLevel, true},
		{"level <= info", common.InfoLevel, true},
		{"level <= info", common.FatalLevel, false},
	}

	var rule common.LogRule
	for _, child := range matrix {
		rule.When = child.code

		if rule.ShouldLog(child.input) == !child.output {
			t.Errorf("%s should return %t with ( %s ), returned %t", child.code, child.output, common.Levels[child.input], !child.output)
		}
	}
}
