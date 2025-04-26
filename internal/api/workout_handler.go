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
		http.Error(w, fmt.Sprintf("failed to parse workoutId %s: %s", workoutIDParam, err.Error()), http.StatusInternalServerError)
		return
	}

	workout, err := wh.workoutStore.FetchWorkout(int(workoutID))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to retrieve workout %d: %s", workoutID, err.Error()), http.StatusInternalServerError)
		return
	}

	if workout == nil {
		http.Error(w, fmt.Sprintf("workout %d not found", workoutID), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(workout)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to write response: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func (wh *WorkoutHandler) CreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout

	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to decode request body: %s", err.Error()), http.StatusBadRequest)
		return
	}

	createdWorkout, err := wh.workoutStore.PersistWorkout(&workout)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to persist workout: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(createdWorkout)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to write response: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func (wh *WorkoutHandler) UpdateWorkout(w http.ResponseWriter, r *http.Request) {
	workoutIDParam := chi.URLParam(r, "workoutId")
	if workoutIDParam == "" {
		http.Error(w, "missing workoutId query param", http.StatusBadRequest)
		return
	}

	workoutID, err := strconv.ParseInt(workoutIDParam, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse workoutId %s: %s", workoutIDParam, err.Error()), http.StatusInternalServerError)
		return
	}

	existingWorkout, err := wh.workoutStore.FetchWorkout(int(workoutID))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to retrieve workout %d: %s", workoutID, err.Error()), http.StatusInternalServerError)
		return
	}

	if existingWorkout == nil {
		http.Error(w, fmt.Sprintf("workout %d not found", workoutID), http.StatusNotFound)
		return
	}

	var updateWorkoutPayload struct {
		Name            *string                 `json:"name"`
		Description     *string                 `json:"description"`
		DurationMinutes *int                    `json:"duration_minutes"`
		CaloriesBurned  *int                    `json:"calories_burned"`
		Exercises       []store.WorkoutExercise `json:"exercises"`
	}

	err = json.NewDecoder(r.Body).Decode(&updateWorkoutPayload)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to decode request body: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if updateWorkoutPayload.Name != nil {
		existingWorkout.Name = *updateWorkoutPayload.Name
	}
	if updateWorkoutPayload.Description != nil {
		existingWorkout.Description = *updateWorkoutPayload.Description
	}
	if updateWorkoutPayload.DurationMinutes != nil {
		existingWorkout.DurationMinutes = *updateWorkoutPayload.DurationMinutes
	}
	if updateWorkoutPayload.CaloriesBurned != nil {
		existingWorkout.CaloriesBurned = *updateWorkoutPayload.CaloriesBurned
	}
	if updateWorkoutPayload.Exercises != nil {
		existingWorkout.Exercises = updateWorkoutPayload.Exercises
	}

	err = wh.workoutStore.UpdateWorkout(existingWorkout)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to update workout: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(existingWorkout)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to write response: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}
