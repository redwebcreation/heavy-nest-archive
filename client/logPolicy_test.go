package client_test

import (
	"testing"

	"github.com/redwebcreation/hez/client"
)

func TestShouldLog(t *testing.T) {
	defaultPolicy := client.LogPolicy{}

	if !(defaultPolicy.ShouldLog(client.DebugLevel) && defaultPolicy.ShouldLog(client.FatalLevel) && defaultPolicy.ShouldLog(client.ErrorLevel)) {
		t.Errorf("an empty policy should always log everything")
	}

	simplePolicy := client.LogPolicy{
		Level: "error",
	}

	if !(simplePolicy.ShouldLog(client.DebugLevel) == false && simplePolicy.ShouldLog(client.ErrorLevel) == true && simplePolicy.ShouldLog(client.FatalLevel) == true) {
		t.Errorf("error in simple policy (weird)")
	}

	matrix := []struct {
		code   string
		input  int
		output bool
	}{
		{"level > error", client.DebugLevel, false},
		{"level > error", client.ErrorLevel, false},
		{"level > error", client.FatalLevel, true},

		{"error > level", client.DebugLevel, true},
		{"error > level", client.ErrorLevel, false},
		{"error > level", client.FatalLevel, false},

		{"level == info", client.InfoLevel, true},
		{"level == info", client.DebugLevel, false},

		{"info == level", client.InfoLevel, true},
		{"info == level", client.DebugLevel, false},

		{"level != info", client.InfoLevel, false},
		{"level != info", client.DebugLevel, true},

		{"info != level", client.InfoLevel, false},
		{"info != level", client.DebugLevel, true},

		{"level < warning", client.DebugLevel, true},
		{"level < warning", client.WarningLevel, false},
		{"level < warning", client.FatalLevel, false},

		{"warning < level", client.DebugLevel, false},
		{"warning < level", client.WarningLevel, false},
		{"warning < level", client.FatalLevel, true},

		{"info >= level", client.DebugLevel, true},
		{"info >= level", client.InfoLevel, true},
		{"info >= level", client.FatalLevel, false},

		{"level >= info", client.DebugLevel, false},
		{"level >= info", client.InfoLevel, true},
		{"level >= info", client.FatalLevel, true},

		{"info <= level", client.DebugLevel, false},
		{"info <= level", client.InfoLevel, true},
		{"info <= level", client.FatalLevel, true},

		{"level <= info", client.DebugLevel, true},
		{"level <= info", client.InfoLevel, true},
		{"level <= info", client.FatalLevel, false},
	}

	var policy client.LogPolicy
	for _, child := range matrix {
		policy.When = child.code

		if policy.ShouldLog(child.input) == !child.output {
			t.Errorf("%s should return %t with ( %s ), returned %t", child.code, child.output, client.Levels[child.input], !child.output)
		}
	}
}
