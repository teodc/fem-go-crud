package routes

import (
	"fem-go-crud/internal/app"
	"github.com/go-chi/chi/v5"
)

func MakeRouter(app *app.App) (r *chi.Mux) {
	r = chi.NewRouter()

	r.Get("/poke", app.HealthCheck)

	r.Get("/users/{userId}", app.UserHandler.GetUser)
	r.Post("/users", app.UserHandler.RegisterUser)

	r.Post("/tokens/authenticate", app.TokenHandler.CreateToken)

	r.Get("/workouts/{workoutId}", app.WorkoutHandler.GetWorkout)
	r.Post("/workouts", app.WorkoutHandler.CreateWorkout)
	r.Put("/workouts/{workoutId}", app.WorkoutHandler.UpdateWorkout)
	r.Delete("/workouts/{workoutId}", app.WorkoutHandler.DeleteWorkout)

	return
}
