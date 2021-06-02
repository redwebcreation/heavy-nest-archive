package core

func ConfigFile() string {
	return ConfigDirectory() + "/hez.yml"
}

func ConfigDirectory() string {
	return "/etc/hez"
}
