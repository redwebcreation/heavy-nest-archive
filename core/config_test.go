package core

import (
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func createFakeConfig(config string) string {
	path := "/tmp/hez_config_fakes/" + strconv.FormatInt(time.Now().Unix(), 10)
	_ = os.RemoveAll(path)

	_ = os.MkdirAll(path, os.FileMode(0777))

	if len(config) == 0 {
		config = `
proxy:
	port: 80
	ssl: 443
	self_signed: false
    logs:
	  level: 0
	  redirections:
		- stdout
applications: []
`
	}

	_ = os.WriteFile(path+"/hez.yml", []byte(strings.TrimSpace(config)), os.FileMode(0777))

	return path + "/hez.yml"
}

func TestFindConfig(t *testing.T) {
	path := createFakeConfig("")

	configFile := FindConfig(path)

	if path != string(configFile) {
		t.Errorf("The config path should be: %s, got %s", path, string(configFile))
	}
}

func TestConfig_IsValid(t *testing.T) {
	path := createFakeConfig("")

	configFile := FindConfig(path)

	if !configFile.IsValid() {
		t.Error("The config file is valid but isValid says it isn't.")
	}
}

func TestConfig_IsValid2(t *testing.T) {
	configFile := FindConfig("/some/random/path/hez.yml")

	if configFile.IsValid() {
		t.Errorf("The config file is invalid but isValid says it is.")
	}
}

func TestConfig_Checksum(t *testing.T) {
	path := createFakeConfig("")

	configFile := FindConfig(path)

	checksum, _ := configFile.Checksum()

	if checksum != "0aba345a632ee0b998958edd6479a69007bcba99d86ff2633ee4b847793ab6c4" {
		t.Errorf("The checksum should be c66f8054c9ffadf3166e694784ddbf22f92c586a4aeb0b8f27a8a666c35a6657, given: %s", checksum)
	}
}

func TestConfig_Resolve(t *testing.T) {
	path := createFakeConfig(`
proxy:
  port: 8080
  ssl: 8443
  self_signed: true
  logs:
    level: 4
    redirections:
      - /tmp/app.log
      - stderr
applications:
  - image: example
    env:
      - APP_ENV=local
    bindings:
      - host: example.com
        container_port: 8000
`)

	configFile := FindConfig(path)

	resolved, err := configFile.Resolve()

	if err != nil {
		t.Error(err)
	}

	if resolved.Applications[0].Image != "example" {
		t.Errorf("applications[0].image should be example, got %s", resolved.Applications[0].Image)
	}

	if resolved.Applications[0].Bindings[0].Host != "example.com" {
		t.Errorf("applications[0].bindings[0].host should be example.com, got %s", resolved.Applications[0].Bindings[0].Host)
	}

	if resolved.Applications[0].Bindings[0].ContainerPort != "8000" {
		t.Errorf("applications[0].bindings[0].container_port should be 8000, got %s", resolved.Applications[0].Bindings[0].ContainerPort)
	}

	if resolved.Applications[0].Env[0] != "APP_ENV=local" {
		t.Errorf("applications[0].env[0] should be APP_ENV=local, got %s", resolved.Applications[0].Env[0])
	}

	if resolved.Proxy.Logs.Level != 4 {
		t.Errorf("proxy.logs.level should be 4, got %d", resolved.Proxy.Logs.Level)
	}

	if resolved.Proxy.Logs.Redirections[0] != "/tmp/app.log" {
		t.Errorf("proxy.logs.redirections[0] should be /tmp/app.log, got %s", resolved.Proxy.Logs.Redirections[0])
	}
	if resolved.Proxy.Logs.Level != 4 {
		t.Errorf("proxy.logs.redirections[1] should be stderr, got %s", resolved.Proxy.Logs.Redirections[1])
	}

	if resolved.Proxy.Port != 8080 {
		t.Errorf("proxy.port should be 8080, got %d", resolved.Proxy.Port)
	}

	if resolved.Proxy.Ssl != 8443 {
		t.Errorf("proxy.ssl should be 8443, got %d", resolved.Proxy.Ssl)
	}

	if !*resolved.Proxy.SelfSigned {
		t.Errorf("proxy.self_signed should be true, got %t", *resolved.Proxy.SelfSigned)
	}

}
