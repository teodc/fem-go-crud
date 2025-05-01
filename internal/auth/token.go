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
	UserID    int64     `json:"-"`
	ExpiresAt time.Time `json:"expires_at"`
	Scope     string    `json:"-"`
}

func MakeToken(userID int64, ttl time.Duration, scope string) (*Token, error) {
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
	hashValue := sha256.Sum256([]byte(token.Plain))
	token.Hash = hashValue[:]

	return token, nil
}
