package common

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ParseMemoryQuota returns the memory size in bytes
func ParseMemoryQuota(raw string) (int64, error) {
	// "unlimited" makes the config more declarative
	if raw == "" || raw == "-1" || raw == "unlimited" {
		return -1, nil
	}

	matches := regexp.MustCompile("^([0-9.]+)([a-z]?)$").FindStringSubmatch(strings.ToLower(raw))

	if len(matches) == 0 {
		return 0, fmt.Errorf("failed to parse memory quantity %q", raw)
	}

	quantity, unit := matches[1], matches[2]

	memory, _ := strconv.ParseFloat(quantity, 64)
	switch unit {
	case "k":
		return int64(memory * 1024.0), nil
	case "":
	case "m":
		return int64(memory * 1024.0 * 1024.0), nil
	case "g":
		return int64(memory * 1024.0 * 1024.0 * 1024.0), nil
	case "t":
		return int64(memory * 1024.0 * 1024.0 * 1024.0 * 1024.0), nil
	}
	return 0, fmt.Errorf("unknown unit %q in memory quantity %q", unit, raw)
}

func MustParseMemoryQuota(raw string) int64 {
	memory, err := ParseMemoryQuota(raw)
	if err != nil {
		panic(err)
	}
	return memory
}
