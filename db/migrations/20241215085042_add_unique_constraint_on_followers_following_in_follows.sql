-- +goose Up
-- +goose StatementBegin
ALTER TABLE follows
ADD CONSTRAINT unique_follow UNIQUE (follower_id, following_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE follows
DROP CONSTRAINT unique_follow;
-- +goose StatementEnd
