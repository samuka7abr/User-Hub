package domain

import (
	"crypto/rand"
	"encoding/base64"
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	PassHash  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

func NewID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
