package store

import (
	"database/sql"
	"errors"
)

type Workout struct {
	ID              int               `json:"id"`
	UserID          int               `json:"user_id"`
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
	PersistWorkout(workout *Workout) error
	GetWorkout(id int) (*Workout, error)
	UpdateWorkout(workout *Workout) error
	DeleteWorkout(id int) error
	GetWorkoutOwner(id int) (int, error)
}

var _ WorkoutStore = (*PostgresWorkoutStore)(nil)

type PostgresWorkoutStore struct {
	db *sql.DB
}

func NewPostgresWorkoutStore(db *sql.DB) *PostgresWorkoutStore {
	return &PostgresWorkoutStore{
		db: db,
	}
}

func (ws *PostgresWorkoutStore) PersistWorkout(workout *Workout) error {
	tx, err := ws.db.Begin()
	if err != nil {
		return err
	}

	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	query := `
		INSERT INTO workouts (user_id, name, description, duration_minutes, calories_burned)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	err = tx.QueryRow(query, workout.UserID, workout.Name, workout.Description, workout.DurationMinutes, workout.CaloriesBurned).Scan(&workout.ID)
	if err != nil {
		return err
	}

	for i := range workout.Exercises {
		// Question: Why do we need that?
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

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (ws *PostgresWorkoutStore) GetWorkout(id int) (*Workout, error) {
	workout := &Workout{}

	workoutQuery := `
		SELECT id, name, description, duration_minutes, calories_burned
		FROM workouts
		WHERE id = $1
	`

	err := ws.db.QueryRow(workoutQuery, id).Scan(
		&workout.ID,
		&workout.Name,
		&workout.Description,
		&workout.DurationMinutes,
		&workout.CaloriesBurned,
	)
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

	rows, err := ws.db.Query(exerciseQuery, id)
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

func (ws *PostgresWorkoutStore) UpdateWorkout(workout *Workout) error {
	tx, err := ws.db.Begin()
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

func (ws *PostgresWorkoutStore) DeleteWorkout(id int) error {
	query := `DELETE FROM workouts WHERE id = $1`

	result, err := ws.db.Exec(query, id)
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

func (ws *PostgresWorkoutStore) GetWorkoutOwner(id int) (int, error) {
	var userID int

	query := `SELECT user_id FROM workouts WHERE id = $1`

	err := ws.db.QueryRow(query, id).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}
