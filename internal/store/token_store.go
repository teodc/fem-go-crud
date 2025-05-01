package store

import (
	"database/sql"

	"fem-go-crud/internal/auth"
)

type TokenStore interface {
	PersistToken(token *auth.Token) error
	RevokeTokensForUser(userID int, scope string) error
}

var _ TokenStore = (*PostgresTokenStore)(nil)

type PostgresTokenStore struct {
	db *sql.DB
}

func NewPostgresTokenStore(db *sql.DB) *PostgresTokenStore {
	return &PostgresTokenStore{
		db: db,
	}
}

func (ts *PostgresTokenStore) PersistToken(token *auth.Token) error {
	query := `
		INSERT INTO tokens (hash, user_id, expires_at, scope)
		VALUES ($1, $2, $3, $4)
	`

	_, err := ts.db.Exec(query, token.Hash, token.UserID, token.ExpiresAt, token.Scope)
	if err != nil {
		return err
	}

	return nil
}

func (ts *PostgresTokenStore) RevokeTokensForUser(userID int, scope string) error {
	query := `
		DELETE FROM tokens
		WHERE user_id = $1 AND scope = $2
	`
	_, err := ts.db.Exec(query, userID, scope)
	if err != nil {
		return err
	}

	return nil
}
