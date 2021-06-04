package core

import (
	"os"
)

func GetKey(key string) (string, error) {
	filePath := StorageDirectory() + "/" + key
	bytes, err := os.ReadFile(filePath)
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
	_ = os.MkdirAll(StorageDirectory(), os.FileMode(0700))

	filePath := StorageDirectory() + "/" + key

	err := os.WriteFile(filePath, []byte(value), os.FileMode(0777))

	if err != nil {
		return err
	}

	return nil
}
