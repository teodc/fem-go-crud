package routes

import (
	"fem-go-crud/internal/app"
	"github.com/go-chi/chi/v5"
)

func MakeRouter(app *app.App) (r *chi.Mux) {
	r = chi.NewRouter()

	r.Get("/poke", app.HealthCheck)
	r.Get("/workouts/{workoutId}", app.WorkoutHandler.GetWorkout)
	r.Post("/workouts", app.WorkoutHandler.CreateWorkout)
	r.Put("/workouts/{workoutId}", app.WorkoutHandler.UpdateWorkout)
	r.Delete("/workouts/{workoutId}", app.WorkoutHandler.DeleteWorkout)

	return
}
