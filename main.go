package main

import (
	"log"
	"net/http"
	"time"

	"fem-go-crud/internal/app"
)

func main() {
	myApp, err := app.New()
	if err != nil {
		panic(err)
	}

	myApp.Logger.Println("Application started")

	server := http.Server{
		Addr:         ":8080",
		Handler:      nil,
		ErrorLog:     myApp.Logger,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  time.Minute,
	}

	http.HandleFunc("/poke", HealthCheck)

	if err := server.ListenAndServe(); err != nil {
		myApp.Logger.Fatal(err)
	}
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("I'm alive!"))
	if err != nil {
		log.Fatal(err)
	}
}
