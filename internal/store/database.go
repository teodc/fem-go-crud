package store

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func Connect() (*sql.DB, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("missing DATABASE_URL env variable")
	}

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	fmt.Println("Database connection opened")

	return db, nil
}

func Migrate(db *sql.DB, migrationsFS fs.FS, migrationsDir string) error {
	goose.SetBaseFS(migrationsFS)
	defer goose.SetBaseFS(nil)

	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("failed to set migrations dialect: %w", err)
	}

	err = goose.Up(db, migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	fmt.Println("Migrations ran successfully")

	return nil
}
