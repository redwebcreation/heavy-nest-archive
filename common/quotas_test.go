package common

import "testing"

func TestMustParseMemoryQuota(t *testing.T) {
	// test it returns one megabyte for 1m
	if mem := MustParseMemoryQuota("1m"); mem != 1024*1024 {
		t.Errorf("Expected 1m to be 1024*1024, got %d", mem)
	}

	// test it returns one gigabyte for 1g
	if mem := MustParseMemoryQuota("1g"); mem != 1024*1024*1024 {
		t.Errorf("Expected 1g to be 1024*1024*1024, got %d", mem)
	}

	// test it returns one gigabyte for 1G
	if mem := MustParseMemoryQuota("1G"); mem != 1024*1024*1024 {
		t.Errorf("Expected 1G to be 1024*1024*1024, got %d", mem)
	}

	// test it returns -1 for an empty string
	if mem := MustParseMemoryQuota(""); mem != -1 {
		t.Errorf("Expected empty string to be -1, got %d", mem)
	}

	// test it returns -1 for unlimited
	if mem := MustParseMemoryQuota("unlimited"); mem != -1 {
		t.Errorf("Expected unlimited to be -1, got %d", mem)
	}

	// test unit-less quantities are treated as megabytes
	if mem := MustParseMemoryQuota("1024"); mem != 1024*1024*1024 {
		t.Errorf("Expected 1024 to be 1024*1024*1024, got %d", mem)
	}

	// test it returns an error for invalid strings
	if _, err := ParseMemoryQuota("invalid_value"); err == nil {
		t.Errorf("Expected invalid_value to be an invalid value")
	}
}

func TestMustParseCpuQuota(t *testing.T) {
	// test it returns 1000000000 for 1
	if cpu := MustParseCpuQuota("1"); cpu != 1000000000 {
		t.Errorf("Expected 1 to be 1000000000, got %d", cpu)
	}

	// test it returns 1500000000 for 1.5
	if cpu := MustParseCpuQuota("1.5"); cpu != 1500000000 {
		t.Errorf("Expected 1.5 to be 1500000000, got %d", cpu)
	}

	// test it returns an error for invalid strings
	if _, err := ParseCpuQuota("invalid_value"); err == nil {
		t.Errorf("Expected invalid_value to be an invalid value")
	}
}
