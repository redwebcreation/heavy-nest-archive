package core

import (
	"os"
)

func GetPathFor(key string) string {
	_ = os.MkdirAll(StorageDirectory(), os.FileMode(0700))

	return StorageDirectory() + "/" + key
}

func CreateFile(name string, value []byte, perm os.FileMode) error {
	err := os.WriteFile(GetPathFor(name), value, perm)

	if err != nil {
		return err
	}

	return nil
}
