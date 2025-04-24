package routes

import (
	"fem-go-crud/internal/app"
	"github.com/go-chi/chi/v5"
)

func MakeRouter(app *app.App) (r *chi.Mux) {
	r = chi.NewRouter()
	r.Get("/poke", app.HealthCheck)
	return
}
