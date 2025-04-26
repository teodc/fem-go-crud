package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"fem-go-crud/internal/store"

	"github.com/go-chi/chi/v5"
)

type WorkoutHandler struct {
	workoutStore store.WorkoutStore
}

func NewWorkoutHandler(ws store.WorkoutStore) *WorkoutHandler {
	return &WorkoutHandler{
		workoutStore: ws,
	}
}

func (wh *WorkoutHandler) GetWorkoutById(w http.ResponseWriter, r *http.Request) {
	workoutIDParam := chi.URLParam(r, "workoutId")
	if workoutIDParam == "" {
		http.Error(w, "missing workoutId query param", http.StatusBadRequest)
		return
	}

	workoutID, err := strconv.ParseInt(workoutIDParam, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get workout %s: %s", workoutIDParam, err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(fmt.Sprintf("workoutId: %d\n", workoutID)))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get workout %d: %s", workoutID, err.Error()), http.StatusInternalServerError)
		return
	}
}

func (wh *WorkoutHandler) CreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout

	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create workout: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	createdWorkout, err := wh.workoutStore.PersistWorkout(&workout)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create workout: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(createdWorkout)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create workout: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}
