package core

import (
	"fmt"
	"os"
)

func GetKey(key string, nullIfEmpty bool) string {
	filePath := StorageDirectory() + "/" + key
	bytes, err := os.ReadFile(filePath)
	data := string(bytes)

	if os.IsNotExist(err) && nullIfEmpty {
		return ""
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return data
}

func SetKey(key string, value string) {
	_ = os.MkdirAll(StorageDirectory(), os.FileMode(0700))

	filePath := StorageDirectory() + "/" + key

	err := os.WriteFile(filePath, []byte(value), os.FileMode(0777))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
