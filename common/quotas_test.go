package common

import "testing"

func TestMustParseMemoryQuota(t *testing.T) {
	// test it returns 512 kilobytes for 512k
	if res := MustParseMemoryQuota("512k"); res != 512*1024 {
		t.Errorf("Expected 512k to return 512*1024, got %d", res)
	}

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

	// test it returns one terabyte for 1t
	if mem := MustParseMemoryQuota("1t"); mem != 1024*1024*1024*1024 {
		t.Errorf("Expected 1t to be 1024*1024*1024*1024, got %d", mem)
	}

	// test it returns one and a half gigabytes for 1.5g
	if mem := MustParseMemoryQuota("1.5g"); mem != 1024*1024*1024*1.5 {
        t.Errorf("Expected 1.5g to be 1024*1024*1024*1.5, got %d", mem)
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
