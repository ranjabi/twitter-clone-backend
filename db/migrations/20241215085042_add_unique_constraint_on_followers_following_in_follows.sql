-- +goose Up
-- +goose StatementBegin
ALTER TABLE follows
ADD CONSTRAINT unique_follow UNIQUE (following_id, followers_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE follows
DROP CONSTRAINT unique_follow;
-- +goose StatementEnd
