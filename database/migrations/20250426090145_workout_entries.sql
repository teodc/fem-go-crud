-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS workout_entries (
    id SERIAL PRIMARY KEY,
    workout_id INTEGER NOT NULL REFERENCES workouts(id) ON DELETE CASCADE,
    exercise_name VARCHAR(100) NOT NULL,
    sets SMALLINT NOT NULL,
    reps SMALLINT DEFAULT NULL,
    duration_seconds SMALLINT DEFAULT NULL,
    weight DECIMAL(5,2) DEFAULT NULL,
    notes TEXT DEFAULT NULL,
    order_index SMALLINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    CONSTRAINT no_reps_and_duration_together CHECK (
        (reps IS NOT NULL OR duration_seconds IS NOT NULL) AND (reps IS NULL OR duration_seconds IS NULL)
    )
 );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE workout_entries;
-- +goose StatementEnd
