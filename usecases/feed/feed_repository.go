package feed

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ranjabi/twitter-clone-backend/types"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	ctx           context.Context
	pgConn        *pgxpool.Pool
	rdConn        *redis.Client
	feedWindowLen int
}

func NewRepository(ctx context.Context, pgConn *pgxpool.Pool, rdConn *redis.Client, feedWindowLen int) Repository {
	return Repository{ctx, pgConn, rdConn, feedWindowLen}
}

func (r *Repository) IsCacheExist(userId int) (bool, error) {
	isExist, err := r.rdConn.Exists(r.ctx, r.GetCacheKey(userId)).Result()
	if err != nil {
		return false, err
	}
	return isExist == 1, nil
}

func (r *Repository) GetCacheKey(userId int) string {
	return fmt.Sprintf("feed:userId:%d", userId)
}

func (r *Repository) GetFeedFromCache(userId int, page int) (*types.Feed, error) {
	offset := (page - 1) * r.feedWindowLen

	tweetsId, err := r.rdConn.LRange(r.ctx, r.GetCacheKey(userId), int64(offset), int64(offset+r.feedWindowLen-1)).Result()
	if err != nil {
		return nil, err
	}

	query := `
		SELECT
			t.id,
			t.content,
			t.like_count,
			t.created_at,
			t.modified_at,
			CASE
				WHEN l.user_id IS NOT NULL THEN TRUE ELSE FALSE
			END AS is_liked,
			u.id,
			u.username,
			u.full_name,
			u.profile_image
		FROM
			tweets t
		INNER JOIN
			users u ON u.id = t.user_id
		LEFT JOIN
			likes l ON l.tweet_id = t.id AND l.user_id = @user_id
		WHERE
			t.id = ANY(@tweetsId)
		ORDER BY
			t.created_at DESC, t.user_id DESC
	`
	args := pgx.NamedArgs{
		"tweetsId": tweetsId,
		"user_id":  userId,
		"limit":    r.feedWindowLen,
		"offset":   offset,
	}
	rows, err := r.pgConn.Query(r.ctx, query, args)
	if err != nil {
		return nil, err
	}

	var tweets []types.Tweet
	for rows.Next() {
		var tweet types.Tweet
		var user types.User

		err := rows.Scan(
			&tweet.Id,
			&tweet.Content,
			&tweet.LikeCount,
			&tweet.CreatedAt,
			&tweet.ModifiedAt,
			&tweet.IsLiked,
			&user.Id,
			&user.Username,
			&user.FullName,
			&user.ProfileImage,
		)
		if err != nil {
			return nil, err
		}

		tweet.User = user
		tweets = append(tweets, tweet)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	var nextPageId *int
	if len(tweets) < r.feedWindowLen {
		nextPageId = nil
	} else {
		nextPage := page + 1
		nextPageId = &nextPage
	}

	feed := types.Feed{
		Tweets:     tweets,
		NextPageId: nextPageId,
	}

	fmt.Printf("Feed:\n%v\n", feed)

	return &feed, nil
}

func (r *Repository) GetFeedTweetsId(userId int) ([]int, error) {
	query := `
		SELECT
			tweet_id,
			user_id
		FROM
			feeds
		WHERE
			user_id = @user_id
	`
	args := pgx.NamedArgs{
		"user_id": userId,
	}
	rows, err := r.pgConn.Query(r.ctx, query, args)
	if err != nil {
		return nil, err
	}

	var tweetsId []int
	for rows.Next() {
		var tweetId int

		if err := rows.Scan(tweetId); err != nil {
			return nil, err
		}

		tweetsId = append(tweetsId, tweetId)
	}

	return tweetsId, nil
}
