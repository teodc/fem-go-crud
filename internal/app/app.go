package app

import (
	"log"
	"net/http"
	"os"

	"fem-go-crud/internal/api"
)

type App struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
}

func New() (*App, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	workoutHandler := api.NewWorkoutHandler()

	app := &App{
		Logger:         logger,
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
