package routes

import (
	"fem-go-crud/internal/app"
	"github.com/go-chi/chi/v5"
)

func MakeRouter(app *app.App) (r *chi.Mux) {
	r = chi.NewRouter()

	// public routes
	r.Get("/poke", app.HealthCheck)
	r.Post("/users", app.UserHandler.RegisterUser)
	r.Post("/tokens/authenticate", app.TokenHandler.CreateToken)

	// protected routes
	r.Group(func(r chi.Router) {
		// Question: Can't we just merge the two middlewares into one? (set user in context + check if valid/authorized)
		r.Use(app.UserMiddleware.Authenticate, app.UserMiddleware.RequireUser)
		r.Get("/users/{userId}", app.UserHandler.GetUser)
		r.Get("/workouts/{workoutId}", app.WorkoutHandler.GetWorkout)
		r.Post("/workouts", app.WorkoutHandler.CreateWorkout)
		r.Put("/workouts/{workoutId}", app.WorkoutHandler.UpdateWorkout)
		r.Delete("/workouts/{workoutId}", app.WorkoutHandler.DeleteWorkout)
	})

	return
}
