package core

import (
	"os"
	"time"
)

func LastChangedTimestamp() string {
	contents, err := os.ReadFile(DataDirectory + "/last_apply_execution")

	if err != nil {
		return time.Now().String()
	}

	return string(contents)
}

func RefreshLastChangedTimestamp() error {
	return os.WriteFile(DataDirectory+"/last_changed_timestamp", []byte(time.Now().String()), os.FileMode(0777))
}
