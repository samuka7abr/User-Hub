package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type JWT struct{ secret []byte }

func NewJWT(secret []byte) *JWT { return &JWT{secret: secret} }

func (j *JWT) Make(sub string, ttl time.Duration) (string, error) {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payload := fmt.Sprintf(`{"sub":"%s","iat":%d,"exp":%d}`, sub, time.Now().Unix(), time.Now().Add(ttl).Unix())
	payloadB64 := base64.RawURLEncoding.EncodeToString([]byte(payload))
	unsigned := header + "." + payloadB64
	mac := hmac.New(sha256.New, j.secret)
	mac.Write([]byte(unsigned))
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return unsigned + "." + sig, nil
}

func (j *JWT) Parse(tok string) (string, error) {
	parts := strings.Split(tok, ".")
	if len(parts) != 3 {
		return "", errors.New("bad token")
	}
	unsigned := parts[0] + "." + parts[1]
	mac := hmac.New(sha256.New, j.secret)
	mac.Write([]byte(unsigned))
	want := mac.Sum(nil)
	got, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil || !hmac.Equal(want, got) {
		return "", errors.New("bad sig")
	}
	pl, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", err
	}
	var m map[string]any
	if err := json.Unmarshal(pl, &m); err != nil {
		return "", err
	}
	exp, ok := m["exp"].(float64)
	if !ok || float64(time.Now().Unix()) > exp {
		return "", errors.New("expired")
	}
	sub, _ := m["sub"].(string)
	if sub == "" {
		return "", errors.New("no sub")
	}
	return sub, nil
}
