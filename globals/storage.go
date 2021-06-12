package globals

import "os"

var DataDirectory string

func init() {
	home, err := os.UserHomeDir()

	if err != nil {
		panic(err)
	}

	DataDirectory = home + "/.config/hez/data/"

	err = os.MkdirAll(DataDirectory, os.FileMode(0777))

	if err != nil {
		panic(err)
	}
}
