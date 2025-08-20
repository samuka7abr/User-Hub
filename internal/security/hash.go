package security

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

type Hasher struct{ pepper []byte }

func NewHasher(pepper []byte) *Hasher { return &Hasher{pepper: pepper} }

func (h *Hasher) Hash(pw string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	sum := sha256.Sum256(append(append(salt, []byte(pw)...), h.pepper...))
	return base64.RawURLEncoding.EncodeToString(salt) + ":" + base64.RawURLEncoding.EncodeToString(sum[:]), nil
}

func (h *Hasher) Verify(stored, pw string) (bool, error) {
	i := -1
	for idx, c := range stored {
		if c == ':' {
			i = idx
			break
		}
	}
	if i < 0 {
		return false, errors.New("bad hash")
	}
	salt, err := base64.RawURLEncoding.DecodeString(stored[:i])
	if err != nil {
		return false, err
	}
	want, err := base64.RawURLEncoding.DecodeString(stored[i+1:])
	if err != nil {
		return false, err
	}
	sum := sha256.Sum256(append(append(salt, []byte(pw)...), h.pepper...))
	return hmac.Equal(sum[:], want), nil
}
