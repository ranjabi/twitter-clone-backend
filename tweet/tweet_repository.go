package tweet

import (
	"context"
	"time"
	"twitter-clone-backend/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	conn *pgxpool.Pool
	ctx  context.Context
}

func NewRepository(conn *pgxpool.Pool, ctx context.Context) Repository {
	return Repository{conn: conn, ctx: ctx}
}

func (r Repository) CreateTweet(tweet model.Tweet) (*model.Tweet, error) {
	var newTweet model.Tweet
	query := `INSERT INTO tweets (content, user_id)  VALUES (@content, @user_id) RETURNING id, content, created_at, user_id`
	args := pgx.NamedArgs{
		"content": tweet.Content,
		"user_id": tweet.UserId,
	}

	err := r.conn.QueryRow(r.ctx, query, args).Scan(&newTweet.Id, &newTweet.Content, &newTweet.CreatedAt, &newTweet.UserId)
	if err != nil {
		return nil, err
	}

	return &newTweet, nil
}

func (r Repository) UpdateTweet(tweet model.Tweet) (*model.Tweet, error) {
	var updatedTweet model.Tweet
	query := `UPDATE tweets SET content=@content, modified_at=@modifiedAt WHERE id=@tweetId RETURNING id, content, modified_at, user_id`
	args := pgx.NamedArgs{
		"tweetId":    tweet.Id,
		"content":    tweet.Content,
		"modifiedAt": time.Now(),
	}

	err := r.conn.QueryRow(r.ctx, query, args).Scan(&updatedTweet.Id, &updatedTweet.Content, &updatedTweet.ModifiedAt, &updatedTweet.UserId)
	if err != nil {
		return nil, err
	}

	return &updatedTweet, nil
}
