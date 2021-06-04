package core

import (
	"os"
	"testing"
)

func TestStorageDirectory(t *testing.T) {
	storageDirectory := StorageDirectory()

	home, _ := os.UserHomeDir()

	if storageDirectory != home+"/.config/hez/storage" {
		t.Errorf("The storage should be: %s, got: %s", home+"/.config/hez/storage", storageDirectory)
	}
}

func TestConfigDirectory(t *testing.T) {
	if ConfigDirectory() != "/etc/hez" {
		t.Errorf("The config directory should be: %s, got: %s", "/etc/hez", ConfigDirectory())
	}
}

func TestConfigFile(t *testing.T) {
	if ConfigFile() != "/etc/hez/hez.yml" {
		t.Errorf("The config file path should be: %s, got: %s", "/etc/hez/hez.yml", ConfigFile())
	}
}
