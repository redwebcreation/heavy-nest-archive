package common_test

import (
	"encoding/base64"
	"encoding/json"
	"github.com/wormable/nest/common"
	"testing"
)

func TestToBase64(t *testing.T) {
	registry := common.RegistryConfiguration{
		Username: "username",
		Password: "password",
	}

	b := registry.ToBase64()

	// decode base64 into text
	payload, err := base64.StdEncoding.DecodeString(b)

	if err != nil {
		t.Errorf(err.Error())
	}

	bytes, _ := json.Marshal(map[string]string{
		"username": "username",
		"password": "password",
	})

	if string(payload) != string(bytes) {
		t.Errorf("Expected %s, got %s", string(bytes), string(payload))
	}
}
