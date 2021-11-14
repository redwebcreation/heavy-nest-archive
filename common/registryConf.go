package common

import (
	"encoding/base64"
	"encoding/json"
	"github.com/wormable/ui"
)

type RegistryConfiguration struct {
	Name string `json:"name"`
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r RegistryConfiguration) ToBase64() string {
	auth, err := json.Marshal(map[string]string{
		"username": r.Username,
		"password": r.Password,
	})
	ui.Check(err)
	return base64.StdEncoding.EncodeToString(auth)
}
