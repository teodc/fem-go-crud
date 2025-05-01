-- +goose Up
-- +goose StatementBegin
ALTER TABLE workouts
ADD COLUMN user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE AFTER
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE workouts
DROP COLUMN user_id;
-- +goose StatementEnd
