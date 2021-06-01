package core

import "strings"

func IsProxyEnabled() bool {
	return false
}

func EnableProxy(port string, ssl string, selfSigned bool) error {
	return nil
}

func DisableProxy() error {
	return nil
}

func GetSupervisordConfig(port string, ssl string, selfSigned bool) string {
	stub := `[program:hezproxy]
directory=/usr/local 
command=/usr/local/bin/hez proxy run --port [port] --ssl [ssl] [selfSigned]
autostart=true
autorestart=true
stderr_logfile=/var/log/hezproxy.err
stdout_logfile=/var/log/hezproxy.log`

	stub = strings.Replace(stub, "[port]", port, 1)
	stub = strings.Replace(stub, "[ssl]", ssl, 1)

	if selfSigned {
		stub = strings.Replace(stub, "[selfSigned]", "--self-signed", 1)
	} else {
		stub = strings.Replace(stub, "[selfSigned]", "", 1)
	}

	return stub
}
