package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/joho/godotenv"

	"fem-go-crud/internal/app"
	"fem-go-crud/internal/routes"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	var port int
	flag.IntVar(&port, "port", 8080, "port to listen on")
	flag.Parse()

	myApp, err := app.New()
	if err != nil {
		panic(err)
	}

	defer func(DB *sql.DB) {
		err := DB.Close()
		if err != nil {
		}
	}(myApp.DB)

	myApp.Logger.Printf("Server started on port %d", port)

	server := http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      routes.MakeRouter(myApp),
		ErrorLog:     myApp.Logger,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  time.Minute,
	}

	if err := server.ListenAndServe(); err != nil {
		myApp.Logger.Fatal(err)
	}
}
