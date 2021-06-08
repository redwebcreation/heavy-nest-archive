package core

import (
	"os"
)

func GetKey(key string) (string, error) {
	bytes, err := os.ReadFile(GetPathFor(key))
	data := string(bytes)

	if os.IsNotExist(err) {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	return data, nil
}

func SetKey(key string, value string) error {
	err := os.WriteFile(GetPathFor(key), []byte(value), os.FileMode(0777))

	if err != nil {
		return err
	}

	return nil
}

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
