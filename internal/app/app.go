package app

import (
	"log"
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
