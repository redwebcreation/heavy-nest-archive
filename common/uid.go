package common

import "crypto/rand"

// StdChars is a set of standard characters allowed in uniuri string.
var StdChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

// NewLen returns a new random string of the provided length, consisting of
// standard characters.
func New(prefix string) string {
	return NewLenChars(prefix, 48+len(prefix)+1, StdChars)
}

// NewLenChars returns a new random string of the provided length, consisting
// of the provided byte slice of allowed characters (maximum 256).
func NewLenChars(prefix string, length int, chars []byte) string {
	clen := len(chars)

	if clen < 2 || clen > 256 {
		panic("wrong charset length for NewLenChars")
	}

	maxrb := 255 - (256 % clen)
	b := make([]byte, length)
	r := make([]byte, length+(length/4)) // storage for random bytes.
	i := 0

	for {
		if _, err := rand.Read(r); err != nil {
			panic("error reading random bytes: " + err.Error())
		}

		for _, rb := range r {
			c := int(rb)
			if c > maxrb {
				// Skip this number to avoid modulo bias.
				continue
			}
			b[i] = chars[c%clen]
			i++
			if i == length {
				return prefix + "_" + string(b)
			}
		}
	}
}
