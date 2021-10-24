package common

import (
	"encoding/base64"
	"encoding/json"

	"github.com/wormable/ui"
)

type RegistryConfiguration struct {
	Host     string
	Username string
	Password string
}

func (r RegistryConfiguration) ToBase64() string {
	auth, err := json.Marshal(map[string]string{
		"username": r.Username,
		"password": r.Password,
	})
	ui.Check(err)
	return base64.StdEncoding.EncodeToString(auth)
}
