package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"
)

const (
	TokenTTL       = 24 * time.Hour
	TokenScopeAuth = "authentication"
)

type Token struct {
	Plain     string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserID    int       `json:"-"`
	ExpiresAt time.Time `json:"expires_at"`
	Scope     string    `json:"-"`
}

func MakeToken(userID int, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID:    userID,
		ExpiresAt: time.Now().Add(ttl),
		Scope:     scope,
	}

	emptyBytes := make([]byte, 32)
	_, err := rand.Read(emptyBytes)
	if err != nil {
		return nil, err
	}

	token.Plain = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(emptyBytes)
	token.Hash = MakeTokenHash(token.Plain)

	return token, nil
}

func MakeTokenHash(plainToken string) []byte {
	hashValue := sha256.Sum256([]byte(plainToken))

	return hashValue[:]
}
