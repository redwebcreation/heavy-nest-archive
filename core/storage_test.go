package core

import (
	"os"
	"testing"
)

func TestSetKey(t *testing.T) {
	_, err := os.Stat(StorageDirectory() + "/my_key")

	if !os.IsNotExist(err) {
		t.Error("The test key `my_key` already exists.")
	}

	setKeyError := SetKey("my_key", "my_value")

	_, err = os.Stat(StorageDirectory() + "/my_key")

	if err != nil {
		t.Errorf("The test key `my_key` does not exists after creation. (%s)", setKeyError)
	}

	_ = os.Remove(StorageDirectory() + "/my_key")
}

func TestGetKey(t *testing.T) {
	value, err := GetKey("does_not_exists")

	if err != nil {
		t.Error(err)
	}

	if value != "" {
		t.Errorf("The test key `does_not_exists` exists.")
	}

	_ = SetKey("does_not_exists", "some_value")

	value, _ = GetKey("does_not_exists")

	if value != "some_value" {
		t.Errorf("The test key `does_not_exists` should be equal to: some_value, given: %s", value)
	}

_ = os.Remove(StorageDirectory() + "/does_not_exists")
}
