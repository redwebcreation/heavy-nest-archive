package tokens

import (
	"bytes"
	"os"
	"strings"

	"github.com/wormable/nest/ansi"
	"github.com/wormable/nest/common"
)

type Token string

var TokenStorage string

func Create() *Token {
	nextToken := Token(common.New("stok"))

	f, err := os.OpenFile(TokenStorage, os.O_APPEND|os.O_WRONLY, 0600)
	ansi.Check(err)

	defer f.Close()

	f.WriteString(string(nextToken) + "\n")

	return &nextToken
}

func All() []Token {
	raw, _ := os.ReadFile(TokenStorage)
	rawTokens := strings.Split(string(raw), "\n")
	var tokens []Token

	for _, token := range rawTokens {
		trimmed := strings.TrimSpace(token)

		if trimmed == "" {
			continue
		}

		tokens = append(tokens, Token(trimmed))
	}

	return tokens
}

func (t *Token) Exists() bool {
	for _, token := range All() {
		if string(token) == string(*t) {
			return true
		}
	}

	return false
}

func Revoke(token string) {
	tokens := All()

	var newTokens bytes.Buffer

	for _, t := range tokens {
		if t != Token(token) {
			newTokens.WriteString(string(t) + "\n")
		}
	}

	os.WriteFile(TokenStorage, newTokens.Bytes(), 0600)
}

func init() {
	storage := common.DataDirectory + "/tokens"

	_, err := os.Stat(storage)

	if os.IsNotExist(err) {
		os.WriteFile(storage, []byte{}, 0600)
	} else {
		ansi.Check(err)
	}

	TokenStorage = storage
}
