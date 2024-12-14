-- +goose Up
-- +goose StatementBegin
ALTER TABLE tweets
ALTER COLUMN created_at SET DEFAULT CURRENT_TIMESTAMP;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tweets
ALTER COLUMN created_at DROP DEFAULT;
-- +goose StatementEnd
