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

	matches := regexp.MustCompile("^([0-9]+)([kmgt]?)$").FindStringSubmatch(strings.ToLower(raw))

	if len(matches) == 0 {
		return 0, fmt.Errorf("failed to parse memory quantity %q", raw)
	}

	quantity, unit := matches[1], matches[2]

	memory, _ := strconv.ParseInt(quantity, 10, 64)
	switch unit {
	case "k":
		return memory * 1024, nil
	case "m":
		return memory * 1024 * 1024, nil
	case "g":
		return memory * 1024 * 1024 * 1024, nil
	case "t":
		return memory * 1024 * 1024 * 1024 * 1024, nil
	}

	// no unit specified, default to megabytes
	return memory * 1024 * 1024, nil
}

func MustParseMemoryQuota(raw string) int64 {
	memory, err := ParseMemoryQuota(raw)
	if err != nil {
		panic(err)
	}
	return memory
}

func ParseCpuQuota(rawQuota string) (int64, error) {
	quota, err := strconv.ParseFloat(rawQuota, 64)
	if err != nil {
		return 0, err
	}

	// quota to nanoseconds
	return int64(quota * 1000000000), nil
}

func MustParseCpuQuota(rawQuota string) int64 {
	quota, err := ParseCpuQuota(rawQuota)
	if err != nil {
		panic(err)
	}
	return quota
}
