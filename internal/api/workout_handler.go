package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"fem-go-crud/internal/store"
	"fem-go-crud/internal/utils"
)

type WorkoutHandler struct {
	workoutStore store.WorkoutStore
	logger       *log.Logger
}

func NewWorkoutHandler(ws store.WorkoutStore, l *log.Logger) *WorkoutHandler {
	return &WorkoutHandler{
		workoutStore: ws,
		logger:       l,
	}
}

func (wh *WorkoutHandler) GetWorkout(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ParseIDParamFromURL(r, "workoutId")
	if err != nil {
		wh.logger.Printf("ERROR: %v", err)
		_ = utils.WriteJSONResponse(w, http.StatusBadRequest, utils.Envelope{"error": "invalid ID param"})
		return
	}

	workout, err := wh.workoutStore.GetWorkout(workoutID)
	if err != nil {
		wh.logger.Printf("ERROR: %v", err)
		_ = utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to retrieve resource"})
		return
	}

	if workout == nil {
		wh.logger.Printf("ERROR: workout %d not found", workoutID)
		_ = utils.WriteJSONResponse(w, http.StatusNotFound, utils.Envelope{"error": "resource not found"})
		return
	}

	_ = utils.WriteJSONResponse(w, http.StatusOK, utils.Envelope{"workout": workout})
}

func (wh *WorkoutHandler) CreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout

	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		wh.logger.Printf("ERROR: %v", err)
		_ = utils.WriteJSONResponse(w, http.StatusBadRequest, utils.Envelope{"error": "invalid payload"})
		return
	}

	err = wh.workoutStore.PersistWorkout(&workout)
	if err != nil {
		wh.logger.Printf("ERROR: %v", err)
		_ = utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to persist resource"})
		return
	}

	_ = utils.WriteJSONResponse(w, http.StatusCreated, utils.Envelope{"workout": workout})
}

func (wh *WorkoutHandler) UpdateWorkout(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ParseIDParamFromURL(r, "workoutId")
	if err != nil {
		wh.logger.Printf("ERROR: %v", err)
		_ = utils.WriteJSONResponse(w, http.StatusBadRequest, utils.Envelope{"error": "invalid ID param"})
		return
	}

	existingWorkout, err := wh.workoutStore.GetWorkout(workoutID)
	if err != nil {
		wh.logger.Printf("ERROR: %v", err)
		_ = utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to retrieve resource"})
		return
	}

	if existingWorkout == nil {
		wh.logger.Printf("ERROR: workout %d not found", workoutID)
		_ = utils.WriteJSONResponse(w, http.StatusNotFound, utils.Envelope{"error": "resource not found"})
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
		wh.logger.Printf("ERROR: %v", err)
		_ = utils.WriteJSONResponse(w, http.StatusBadRequest, utils.Envelope{"error": "invalid payload"})
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
		wh.logger.Printf("ERROR: %v", err)
		_ = utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to update resource"})
		return
	}

	_ = utils.WriteJSONResponse(w, http.StatusOK, utils.Envelope{"workout": existingWorkout})
}

func (wh *WorkoutHandler) DeleteWorkout(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ParseIDParamFromURL(r, "workoutId")
	if err != nil {
		wh.logger.Printf("ERROR: %v", err)
		_ = utils.WriteJSONResponse(w, http.StatusBadRequest, utils.Envelope{"error": "invalid ID param"})
		return
	}

	err = wh.workoutStore.DeleteWorkout(workoutID)
	// idempotency?
	if errors.Is(err, sql.ErrNoRows) {
		wh.logger.Printf("ERROR: %v", err)
		_ = utils.WriteJSONResponse(w, http.StatusNotFound, utils.Envelope{"error": "resource not found"})
		return
	}
	if err != nil {
		wh.logger.Printf("ERROR: %v", err)
		_ = utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to delete resource"})
		return
	}

	_ = utils.WriteJSONResponse(w, http.StatusNoContent, utils.Envelope{})
}
