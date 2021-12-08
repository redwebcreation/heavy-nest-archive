package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"golang.org/x/term"
	"io"
	"os"
	"regexp"
	"syscall"
)

func AskPrivately(label string) ([]byte, error) {
	fmt.Printf(label)
	password, err := term.ReadPassword(syscall.Stdin)
	fmt.Println("********")
	return password, err
}

type Vault struct {
	Path string
}

type Secret struct {
	Key      string
	Value    []byte
	Password []byte
}

func (v Vault) Put(secret Secret) error {
	key := sha256.Sum256(secret.Password)
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	contents := fmt.Sprintf(
		"$NEST_VAULT;1.0;AES256;key=%s\n%s",
		secret.Key,
		InsertNewLines(aesGCM.Seal(nonce, nonce, secret.Value, nil)),
	)

	return os.WriteFile(v.Path+"/"+secret.Key, []byte(contents), 0600)
}

func InsertNewLines(b []byte) string {
	var r = regexp.MustCompile("(.{64})+")
	return r.ReplaceAllString(fmt.Sprintf("%x", b), "$1\n")
}
