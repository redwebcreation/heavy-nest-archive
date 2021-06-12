package core

import (
	"github.com/redwebcreation/hez2/globals"
	"os"
	"time"
)

func GetLastApplyExecution() string {
	contents, err := os.ReadFile(globals.DataDirectory + "/last_apply_execution")

	if err != nil {
		return time.Now().String()
	}

	return string(contents)
}

func RefreshLastApplyExecution() error {
	return os.WriteFile(globals.DataDirectory+"/last_apply_execution", []byte(time.Now().String()), os.FileMode(0777))
}
