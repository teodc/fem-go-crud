package store

import (
	"database/sql"
	"errors"
)

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
	GetWorkout(id int64) (*Workout, error)
	UpdateWorkout(workout *Workout) error
	DeleteWorkout(id int64) error
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

	for i := range workout.Exercises {
		// Why?
		exercise := &workout.Exercises[i]

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

func (pws *PostgresWorkoutStore) GetWorkout(id int64) (*Workout, error) {
	workout := &Workout{}

	workoutQuery := `
		SELECT id, name, description, duration_minutes, calories_burned
		FROM workouts
		WHERE id = $1
	`

	err := pws.db.QueryRow(workoutQuery, id).Scan(&workout.ID, &workout.Name, &workout.Description, &workout.DurationMinutes, &workout.CaloriesBurned)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	exerciseQuery := `
		SELECT id, name, sets, reps, duration_seconds, weight, notes, order_index
		FROM workout_exercises
		WHERE workout_id = $1
		ORDER BY order_index
	`

	rows, err := pws.db.Query(exerciseQuery, id)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		var exercise WorkoutExercise
		err = rows.Scan(&exercise.ID, &exercise.Name, &exercise.Sets, &exercise.Reps, &exercise.DurationSeconds, &exercise.Weight, &exercise.Notes, &exercise.OrderIndex)
		if err != nil {
			return nil, err
		}

		workout.Exercises = append(workout.Exercises, exercise)
	}

	return workout, nil
}

func (pws *PostgresWorkoutStore) UpdateWorkout(workout *Workout) error {
	tx, err := pws.db.Begin()
	if err != nil {
		return err
	}

	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	query := `
		UPDATE workouts
		SET name = $1, description = $2, duration_minutes = $3, calories_burned = $4
		WHERE id = $5
	`

	result, err := tx.Exec(query, workout.Name, workout.Description, workout.DurationMinutes, workout.CaloriesBurned, workout.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	_, err = tx.Exec(`DELETE FROM workout_exercises WHERE workout_id = $1`, workout.ID)
	if err != nil {
		return err
	}

	for i := range workout.Exercises {
		exercise := &workout.Exercises[i]

		query := `
			INSERT INTO workout_exercises (workout_id, name, sets, reps, duration_seconds, weight, notes, order_index)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id
		`

		err = tx.QueryRow(query, workout.ID, exercise.Name, exercise.Sets, exercise.Reps, exercise.DurationSeconds, exercise.Weight, exercise.Notes, exercise.OrderIndex).Scan(&exercise.ID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (pws *PostgresWorkoutStore) DeleteWorkout(id int64) error {
	query := `DELETE FROM workouts WHERE id = $1`

	result, err := pws.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
