package tweet

import (
	"context"
	"time"
	"twitter-clone-backend/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type TweetRepository struct {
	ctx    context.Context
	pgConn *pgxpool.Pool
}

func NewRepository(ctx context.Context, pgConn *pgxpool.Pool, rdConn *redis.Client) TweetRepository {
	return TweetRepository{ctx: ctx, pgConn: pgConn}
}

func (r *TweetRepository) CreateTweet(tweet models.Tweet) (*models.Tweet, error) {
	var newTweet models.Tweet
	query := `INSERT INTO tweets (content, user_id)  VALUES (@content, @user_id) RETURNING id, content, created_at, user_id`
	args := pgx.NamedArgs{
		"content": tweet.Content,
		"user_id": tweet.UserId,
	}

	err := r.pgConn.QueryRow(r.ctx, query, args).Scan(&newTweet.Id, &newTweet.Content, &newTweet.CreatedAt, &newTweet.UserId)
	if err != nil {
		return nil, err
	}

	return &newTweet, nil
}

func (r *TweetRepository) FindById(id int) (*models.Tweet, error) {
	var tweet models.Tweet
	query := `SELECT id, content, created_at, user_id from tweets WHERE id = @id`
	args := pgx.NamedArgs{
		"id": id,
	}

	err := r.pgConn.QueryRow(r.ctx, query, args).Scan(&tweet.Id, &tweet.Content, &tweet.CreatedAt, &tweet.UserId)
	if err != nil {
		return nil, err
	}

	return &tweet, nil
}

func (r *TweetRepository) IsTweetExistById(id int) (bool, error) {
	var isTweetExist bool
	query := `SELECT EXISTS (SELECT 1 FROM tweets WHERE id=@id)`
	args := pgx.NamedArgs{
		"id": id,
	}

	err := r.pgConn.QueryRow(r.ctx, query, args).Scan(&isTweetExist)
	if err != nil {
		return false, err
	}

	return isTweetExist, nil
}

func (r *TweetRepository) UpdateTweet(tweet models.Tweet) (*models.Tweet, error) {
	var updatedTweet models.Tweet
	query := `UPDATE tweets SET content=@content, modified_at=@modifiedAt WHERE id=@tweetId RETURNING id, content, modified_at, user_id`
	args := pgx.NamedArgs{
		"tweetId":    tweet.Id,
		"content":    tweet.Content,
		"modifiedAt": time.Now(),
	}

	err := r.pgConn.QueryRow(r.ctx, query, args).Scan(&updatedTweet.Id, &updatedTweet.Content, &updatedTweet.ModifiedAt, &updatedTweet.UserId)
	if err != nil {
		return nil, err
	}

	return &updatedTweet, nil
}

func (r *TweetRepository) DeleteTweet(id int) error {
	query := `DELETE FROM tweets WHERE id=@id`
	args := pgx.NamedArgs{
		"id": id,
	}
	if _, err := r.pgConn.Exec(r.ctx, query, args); err != nil {
		return err
	}

	return nil
}

func (r *TweetRepository) IsTweetLiked(userId int, tweetId int) (bool, error) {
	var isTweetLiked bool
	query := `SELECT EXISTS (SELECT 1 FROM likes WHERE user_id=@user_id AND tweet_id=@tweet_id)`
	args := pgx.NamedArgs{
		"user_id":  userId,
		"tweet_id": tweetId,
	}

	err := r.pgConn.QueryRow(r.ctx, query, args).Scan(&isTweetLiked)
	if err != nil {
		return false, err
	}

	return isTweetLiked, nil
}

func (r *TweetRepository) GetTweetLikeCountById(id int) (int, error) {
	var likeCount int
	query := `SELECT like_count from tweets WHERE id = @id`
	args := pgx.NamedArgs{
		"id": id,
	}

	err := r.pgConn.QueryRow(r.ctx, query, args).Scan(&likeCount)
	if err != nil {
		return 0, err
	}

	return likeCount, nil
}

// TODO: should separate to LikeRepository?
func (r *TweetRepository) LikeTweet(userId int, tweetId int) (int, error) {
	var likeCount int
	tx, err := r.pgConn.Begin(r.ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(r.ctx)

	query := "INSERT INTO likes (user_id, tweet_id) VALUES (@user_id, @tweet_id)"
	args := pgx.NamedArgs{
		"user_id":  userId,
		"tweet_id": tweetId,
	}
	_, err = tx.Exec(r.ctx, query, args)
	if err != nil {
		return 0, err
	}

	query = "UPDATE tweets SET like_count = like_count + 1 WHERE id = @id RETURNING like_count"
	args = pgx.NamedArgs{
		"id": tweetId,
	}
	err = r.pgConn.QueryRow(r.ctx, query, args).Scan(&likeCount)
	if err != nil {
		return 0, err
	}

	err = tx.Commit(r.ctx)
	if err != nil {
		return 0, nil
	}

	return likeCount, nil
}

func (r *TweetRepository) UnlikeTweet(userId int, tweetId int) (int, error) {
	var likeCount int
	tx, err := r.pgConn.Begin(r.ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(r.ctx)

	query := "DELETE FROM likes WHERE user_id = @user_id and tweet_id = @tweet_id"
	args := pgx.NamedArgs{
		"user_id":  userId,
		"tweet_id": tweetId,
	}
	_, err = tx.Exec(r.ctx, query, args)
	if err != nil {
		return 0, err
	}

	query = "UPDATE tweets SET like_count = like_count - 1 WHERE id = @id RETURNING like_count"
	args = pgx.NamedArgs{
		"id": tweetId,
	}
	err = r.pgConn.QueryRow(r.ctx, query, args).Scan(&likeCount)
	if err != nil {
		return 0, err
	}

	err = tx.Commit(r.ctx)
	if err != nil {
		return 0, nil
	}

	return likeCount, nil
}
