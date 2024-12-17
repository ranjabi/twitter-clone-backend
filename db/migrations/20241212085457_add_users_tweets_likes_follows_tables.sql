-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
	id 				int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	username		varchar(16) NOT NULL,
	email 			varchar(255) NOT NULL,
	password 		varchar(128) NOT NULL,
	created_at 		timestamp NOT NULL,
	following_count int DEFAULT 0,
	followers_count	int DEFAULT 0
);

CREATE TABLE tweets (
	id 			int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	content		text NOT NULL,
	created_at 	timestamp NOT NULL,
	modified_at	timestamp,
	user_id 	int NOT NULL REFERENCES users(id)
);

CREATE TABLE likes (
	id 			int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	user_id		int NOT NULL REFERENCES users(id),
	tweet_id	int NOT NULL REFERENCES tweets(id)
);

CREATE TABLE follows (
	id				int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	follower_id	int NOT NULL REFERENCES users(id),
	following_id	int NOT NULL REFERENCES users(id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE follows;
DROP TABLE likes;
DROP TABLE tweets;
DROP TABLE users;
-- +goose StatementEnd
