-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD COLUMN full_name varchar(255) NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
    DROP COLUMN full_name;
-- +goose StatementEnd
