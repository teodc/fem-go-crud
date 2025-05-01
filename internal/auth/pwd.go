package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	Plain *string
	Hash  []byte
}

func (p *Password) Set(plain string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.MinCost)
	if err != nil {
		return err
	}

	p.Plain = &plain
	p.Hash = hash

	return nil
}

func (p *Password) Matches(plain string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(plain))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
