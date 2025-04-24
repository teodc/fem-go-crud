package app

import (
	"log"
	"net/http"
	"os"
)

type App struct {
	Logger *log.Logger
}

func New() (*App, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	app := &App{
		Logger: logger,
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
