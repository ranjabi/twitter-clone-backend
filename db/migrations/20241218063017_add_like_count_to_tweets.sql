-- +goose Up
-- +goose StatementBegin
ALTER TABLE tweets
    ADD COLUMN like_count int DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tweets
    DROP COLUMN like_count;
-- +goose StatementEnd
