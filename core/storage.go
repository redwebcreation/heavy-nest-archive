package core

import (
	"io/ioutil"
	"log"
	"os"
)

func StorageDirectory() string {
	path := configDirectory() + "/compiled"

	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		os.Mkdir(path, 777)
	}

	return path
}

func GetKey(key string) string {
	filePath := StorageDirectory() + "/" + key
	bytes, err := ioutil.ReadFile(filePath)
	data := string(bytes)

	if err != nil {
		log.Fatal(err)
	}

	return data
}

func SetKey(key string, value string) {
	filePath := StorageDirectory() + "/" + key

	_, err := os.Stat(filePath)

	if !os.IsNotExist(err) {
		log.Fatal("Calling SetKey on an already set value is prohibited. Use SetKeyOverride.")
	}

	SetKeyOverride(key, value)
}

func SetKeyOverride(key string, value string) {
	filePath := StorageDirectory() + "/" + key

	file, err := os.Create(filePath)

	if err != nil {
		log.Fatal(err)
	}

	_, err = file.WriteString(value)

	if err != nil {
		file.Close()
		log.Fatal(err)
	}

	err = file.Close()

	if err != nil {
		log.Fatal(err)
	}
}
