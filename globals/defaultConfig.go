package globals

import (
	_ "embed"
)

//go:embed _default.json
var DefaultConfig []byte
