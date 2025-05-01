package store

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"fem-go-crud/internal/auth"
)

type User struct {
	ID        int           `json:"id"`
	Username  string        `json:"username"`
	Email     string        `json:"email"`
	Password  auth.Password `json:"-"`
	CreatedAt string        `json:"created_at"`
	UpdatedAt string        `json:"updated_at"`
}

var AnonymousUser = &User{}

func (u *User) IsAnonymous() bool {
	// Question: Is this comparing 2 pointer addresses? Or both structs' values?
	return u == AnonymousUser
}

type UserStore interface {
	PersistUser(user *User) error
	GetUserByIdOrUsername(id int, username string) (*User, error)
	UpdateUser(user *User) error
	GetUserFromToken(token, scope string) (*User, error)
}

var _ UserStore = (*PostgresUserStore)(nil)

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{
		db: db,
	}
}

func (us *PostgresUserStore) PersistUser(user *User) error {
	query := `
		INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	err := us.db.QueryRow(query, user.Username, user.Email, user.Password.Hash).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (us *PostgresUserStore) GetUserByIdOrUsername(id int, username string) (*User, error) {
	var targetField string
	var targetValue any

	if id != 0 {
		targetField = "id"
		targetValue = id
	} else if username != "" {
		targetField = "username"
		targetValue = username
	} else {
		return nil, errors.New("missing id or username")
	}

	user := &User{
		Password: auth.Password{},
	}

	query := `
		SELECT id, username, email, password_hash, created_at, updated_at
		FROM users
		WHERE %s = $1
	`

	err := us.db.QueryRow(fmt.Sprintf(query, targetField), fmt.Sprintf("%v", targetValue)).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password.Hash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *PostgresUserStore) UpdateUser(user *User) error {
	query := `
		UPDATE users
		SET username = $1, email = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
		RETURNING updated_at
	`

	result, err := us.db.Exec(query, user.Username, user.Email, user.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (us *PostgresUserStore) GetUserFromToken(plainToken, scope string) (*User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.created_at, u.updated_at
		FROM users u
		INNER JOIN tokens t ON u.id = t.user_id
		WHERE t.hash = $1 AND t.scope = $2 AND t.expires_at > $3
	`

	tokenHash := auth.MakeTokenHash(plainToken)

	user := &User{
		Password: auth.Password{},
	}

	err := us.db.QueryRow(query, tokenHash, scope, time.Now()).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}
