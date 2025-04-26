package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type WorkoutHandler struct{}

func NewWorkoutHandler() *WorkoutHandler {
	return &WorkoutHandler{}
}

func (wh *WorkoutHandler) GetWorkoutById(w http.ResponseWriter, r *http.Request) {
	workoutIDParam := chi.URLParam(r, "workoutId")
	if workoutIDParam == "" {
		http.Error(w, "missing workoutId query param", http.StatusBadRequest)
		return
	}

	workoutID, err := strconv.ParseInt(workoutIDParam, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(fmt.Sprintf("workoutId: %d\n", workoutID)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (wh *WorkoutHandler) CreateWorkout(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("workout created"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
