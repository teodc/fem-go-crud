package app

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"fem-go-crud/database/migrations"
	"fem-go-crud/internal/api"
	"fem-go-crud/internal/store"
)

type App struct {
	Logger         *log.Logger
	DB             *sql.DB
	WorkoutHandler *api.WorkoutHandler
}

func New() (*App, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := store.Connect()
	if err != nil {
		return nil, err
	}

	err = store.Migrate(db, migrations.FS, ".")
	if err != nil {
		return nil, err
	}

	workoutHandler := api.NewWorkoutHandler()

	app := &App{
		Logger:         logger,
		DB:             db,
		WorkoutHandler: workoutHandler,
	}

	return app, nil
}

func (a *App) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("I'm alive!"))
	if err != nil {
		log.Fatal(err)
	}
}
