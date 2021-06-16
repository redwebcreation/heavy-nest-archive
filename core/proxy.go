package core

func GetWhitelistedDomains() []string {
	var domains = make([]string, len(Config.Applications))

	for _, application := range Config.Applications {
		domains = append(domains, application.Host)
	}

	return domains
}