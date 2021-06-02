package core

import (
	"fmt"
	"os"
)

func StorageDirectory() string {
	path := ConfigDirectory() + "/compiled"

	EnsureDirectoryExists(path)

	return path
}

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
	filePath := StorageDirectory() + "/" + key

	_, err := os.Stat(filePath)

	if !os.IsNotExist(err) {
		fmt.Println("Calling SetKey on an already set value is prohibited. Use SetKeyOverride.")
		os.Exit(1)
	}

	SetKeyOverride(key, value)
}

func SetKeyOverride(key string, value string) {
	filePath := StorageDirectory() + "/" + key

	file, err := os.Create(filePath)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_, err = file.WriteString(value)

	if err != nil {
		file.Close()
		fmt.Println(err)
		os.Exit(1)
	}

	err = file.Close()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func EnsureDirectoryExists(path string) {
	_, err := os.Stat(path)

	if !os.IsNotExist(err) {
		os.Mkdir(path, os.FileMode(0744))
	}
}
