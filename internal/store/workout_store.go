package store

import "database/sql"

type Workout struct {
	ID              int               `json:"id"`
	Name            string            `json:"name"`
	Description     string            `json:"description"`
	DurationMinutes int               `json:"duration_minutes"`
	CaloriesBurned  int               `json:"calories_burned"`
	Exercises       []WorkoutExercise `json:"exercises"`
}

type WorkoutExercise struct {
	ID              int      `json:"id"`
	Name            string   `json:"name"`
	Sets            int      `json:"sets"`
	Reps            *int     `json:"reps"`
	DurationSeconds *int     `json:"duration_seconds"`
	Weight          *float64 `json:"weight"`
	Notes           string   `json:"notes"`
	OrderIndex      int      `json:"order_index"`
}

type WorkoutStore interface {
	PersistWorkout(workout *Workout) (*Workout, error)
	FetchWorkoutById(id int) (*Workout, error)
}

type PostgresWorkoutStore struct {
	db *sql.DB
}

func NewPostgresWorkoutStore(db *sql.DB) *PostgresWorkoutStore {
	return &PostgresWorkoutStore{
		db: db,
	}
}

func (pws *PostgresWorkoutStore) PersistWorkout(workout *Workout) (*Workout, error) {
	tx, err := pws.db.Begin()
	if err != nil {
		return nil, err
	}

	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	query := `
		INSERT INTO workouts (name, description, duration_minutes, calories_burned)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err = tx.QueryRow(query, workout.Name, workout.Description, workout.DurationMinutes, workout.CaloriesBurned).Scan(&workout.ID)
	if err != nil {
		return nil, err
	}

	for _, exercise := range workout.Exercises {
		query := `
			INSERT INTO workout_exercises (workout_id, name, sets, reps, duration_seconds, weight, notes, order_index)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id
		`

		err = tx.QueryRow(query, workout.ID, exercise.Name, exercise.Sets, exercise.Reps, exercise.DurationSeconds, exercise.Weight, exercise.Notes, exercise.OrderIndex).Scan(&exercise.ID)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return workout, nil
}

func (pws *PostgresWorkoutStore) FetchWorkoutById(id int) (*Workout, error) {
	return &Workout{}, nil
}
