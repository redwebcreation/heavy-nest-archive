package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
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

func (v Vault) Has(key string) bool {
	_, err := os.Stat(v.Path + "/" + key)
	return err == nil
}

func DerivePassword(password []byte) (cipher.AEAD, error) {
	key := sha256.Sum256(password)
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return aesGCM, nil
}

func (v Vault) Put(secret Secret) error {
	aesGCM, err := DerivePassword(secret.Password)
	if err != nil {
		return err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	contents := fmt.Sprintf(
		"$NEST_VAULT;1.0;AES256\n%x",
		aesGCM.Seal(nonce, nonce, secret.Value, nil),
	)

	return os.WriteFile(v.Path+"/"+secret.Key, []byte(contents), 0600)
}

func (v Vault) Get(s *Secret) ([]byte, error) {
	bytes, err := os.ReadFile(v.Path + "/" + s.Key)
	if err != nil {
		return nil, err
	}

	parsed, err := ParseSecret(bytes)
	if err != nil {
		return nil, err
	}

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(s.Password)
	if err != nil {
		return nil, err
	}

	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	//Get the nonce size
	nonceSize := aesGCM.NonceSize()

	//Extract the nonce from the encrypted data
	nonce, ciphertext := parsed.Value[:nonceSize], parsed.Value[nonceSize:]

	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%s\n", plaintext)
	return plaintext, nil
}

type ParsedSecret struct {
	Hash         string
	VaultVersion string
	Value        []byte
}

func ParseSecret(s []byte) (ParsedSecret, error) {
	re := regexp.MustCompile(`(?m)^\$NEST_VAULT;(?P<version>[0-9.]+);(?P<hash>[a-zA-Z1-9]+)\n(?P<contents>[0-9a-f]+)$`)
	matches := re.FindAllSubmatch(s, -1)

	if len(matches) != 1 {
		return ParsedSecret{}, fmt.Errorf("invalid vault file")
	}

	vaultVersion, hash, contents := string(matches[0][1]), string(matches[0][2]), matches[0][3]

	if vaultVersion != "1.0" {
		return ParsedSecret{}, fmt.Errorf("unsupported vault version (supported: 1.0)")
	}

	if hash != "AES256" {
		return ParsedSecret{}, fmt.Errorf("unsupported hash algorithm (supported: AES256)")
	}

	decoded, err := hex.DecodeString(string(contents))
	if err != nil {
		return ParsedSecret{}, err
	}
	return ParsedSecret{
		Hash:         hash,
		VaultVersion: vaultVersion,
		Value:        decoded,
	}, nil
}
