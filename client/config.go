package client

type Config struct {
	DefaultNetwork string

	Applications map[string]struct {
		Image    string
		Env      map[string]string
		EnvFiles []string
		Warm     bool
		Registry RegistryAuth
	}

	Registries map[string]RegistryAuth

	Staging struct {
		Enabled   bool
		Host      string
		LogPolicy LogPolicy

		MaxVersions int // -1 for every commit, n for last n commits available in stating

		Database struct {
			Internal bool

			Type string
			DSN  string
		}

		Applications []string
	}

	Production struct {
		LogPolicy LogPolicy
		HttpPort  string
		HttpsPort string
	}

	LogPolicies map[string]LogPolicy
}

type LogPolicy struct {
	Level        int
	Redirections []string
}
