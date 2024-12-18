-- +goose Up
-- +goose StatementBegin
ALTER TABLE likes
    ADD COLUMN created_at timestamp DEFAULT CURRENT_TIMESTAMP;

ALTER TABLE likes
    DROP CONSTRAINT IF EXISTS likes_user_id_fkey,
    ADD CONSTRAINT likes_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    DROP CONSTRAINT IF EXISTS likes_tweet_id_fkey,
    ADD CONSTRAINT likes_tweet_id_fkey FOREIGN KEY (tweet_id) REFERENCES tweets(id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE likes
    DROP CONSTRAINT likes_user_id_fkey,
    DROP CONSTRAINT likes_tweet_id_fkey;

ALTER TABLE likes
    DROP COLUMN created_at;
-- +goose StatementEnd
