package store

import (
	"database/sql"
	"errors"

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

type UserStore interface {
	PersistUser(user *User) error
	GetUser(id int64) (*User, error)
	UpdateUser(user *User) error
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

func (us *PostgresUserStore) GetUser(id int64) (*User, error) {
	user := &User{
		Password: auth.Password{},
	}

	query := `
		SELECT id, username, email, password_hash, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	err := us.db.QueryRow(query, id).Scan(
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
