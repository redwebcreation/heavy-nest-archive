package vault

import "testing"

func TestParseSecret(t *testing.T) {
	contents := "$NEST_VAULT;1.0;AES256;key=key\na13a6595501fd88f16c65632afa23efe0711ac0a830e781bf302a6a5de811e17\n204c6b98efb7b5"

	ParseSecret([]byte(contents))
}
