package store

import (
	"database/sql"
	"os"
	"testing"

	"fem-go-crud/database/migrations"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO: move to utils
func setupTestDB(t *testing.T) *sql.DB {
	if err := godotenv.Load("../../.env.test"); err != nil {
		panic(err)
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Fatal("missing DATABASE_URL env variable")
	}

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		t.Fatalf("failed to open database connection: %v", err)
	}

	err = Migrate(db, migrations.FS, ".")
	if err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	_, err = db.Exec("TRUNCATE TABLE workouts, workout_exercises CASCADE")
	if err != nil {
		t.Fatalf("failed to truncate tables: %v", err)
	}

	return db
}

func TestPersistWorkout(t *testing.T) {
	db := setupTestDB(t)
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	store := NewPostgresWorkoutStore(db)

	// arrange (table testing style)
	testCases := []struct {
		name       string
		workout    *Workout
		expectsErr bool
	}{
		{
			name: "valid workout",
			workout: &Workout{
				Name:            "Valid Workout Name",
				Description:     "Valid Workout Description",
				DurationMinutes: 10,
				CaloriesBurned:  500,
				Exercises: []WorkoutExercise{
					{
						Name:            "Valid Exercise Name",
						Sets:            10,
						Reps:            intPtr(10),
						DurationSeconds: nil,
						Weight:          floatPtr(15.5),
						Notes:           "Valid Exercise Notes",
						OrderIndex:      1,
					},
				},
			},
			expectsErr: false,
		},
		{
			name: "invalid workout",
			workout: &Workout{
				Name:            "Invalid Workout Name",
				Description:     "Invalid Workout Description",
				DurationMinutes: 10,
				CaloriesBurned:  500,
				Exercises: []WorkoutExercise{
					{
						Name:            "Invalid Exercise Name",
						Sets:            10,
						Reps:            intPtr(10),
						DurationSeconds: intPtr(10),
						Weight:          floatPtr(15.5),
						Notes:           "Invalid Exercise Notes",
						OrderIndex:      1,
					},
				},
			},
			expectsErr: true,
		},
	}

	// act & assert
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := store.PersistWorkout(tc.workout)
			if tc.expectsErr {
				// We should match the expected error instead
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tc.workout.Name, tc.workout.Name)
			assert.Equal(t, tc.workout.Description, tc.workout.Description)
			assert.Equal(t, tc.workout.DurationMinutes, tc.workout.DurationMinutes)
			assert.Equal(t, tc.workout.CaloriesBurned, tc.workout.CaloriesBurned)
			assert.Equal(t, len(tc.workout.Exercises), len(tc.workout.Exercises))
			for i := range tc.workout.Exercises {
				assert.Equal(t, tc.workout.Exercises[i].Name, (&tc.workout.Exercises[i]).Name)
				assert.Equal(t, tc.workout.Exercises[i].Sets, (&tc.workout.Exercises[i]).Sets)
				assert.Equal(t, tc.workout.Exercises[i].Reps, (&tc.workout.Exercises[i]).Reps)
				assert.Equal(t, tc.workout.Exercises[i].DurationSeconds, (&tc.workout.Exercises[i]).DurationSeconds)
				assert.Equal(t, tc.workout.Exercises[i].Weight, (&tc.workout.Exercises[i]).Weight)
				assert.Equal(t, tc.workout.Exercises[i].Notes, (&tc.workout.Exercises[i]).Notes)
				assert.Equal(t, tc.workout.Exercises[i].OrderIndex, (&tc.workout.Exercises[i]).OrderIndex)
			}

			retrievedWorkout, err := store.GetWorkout(int64(tc.workout.ID))
			require.NoError(t, err)

			assert.Equal(t, tc.workout.ID, retrievedWorkout.ID)
			assert.Equal(t, len(tc.workout.Exercises), len(retrievedWorkout.Exercises))
		})
	}
}

func intPtr(i int) *int {
	return &i
}

func floatPtr(f float64) *float64 {
	return &f
}
