package user

import (
	"context"
	"fmt"
	"time"

	"github.com/ranjabi/twitter-clone-backend/types"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

const (
	userProfilePath             = "$"
	userProfileRecentTWeetsPath = "$"
	userProfileExpirationTime   = 10 * time.Minute
)

func getUserProfileCacheKey(id int) string {
	return fmt.Sprintf("user.id:%d", id)
}

func getRecentTweetsCacheKey(id int) string {
	return fmt.Sprintf("user.id:%d:recentTweets", id)
}

type Repository struct {
	ctx    context.Context
	pgConn *pgxpool.Pool
	rdConn *redis.Client
}

func NewRepository(ctx context.Context, pgConn *pgxpool.Pool, rdConn *redis.Client) Repository {
	return Repository{ctx, pgConn, rdConn}
}

func (r *Repository) Create(user types.User) (*types.User, error) {
	query := `
		INSERT INTO users (
			full_name, 
			username, 
			email, 
			password
		) 
		VALUES (
			@full_name, 
			LOWER(@username), 
			LOWER(@email), 
			@password
		) 
		RETURNING 
			id, 
			full_name, 
			username, 
			email
	`
	args := pgx.NamedArgs{
		"full_name": user.FullName,
		"username":  user.Username,
		"email":     user.Email,
		"password":  user.Password,
	}

	var newUser types.User
	err := r.pgConn.QueryRow(r.ctx, query, args).Scan(
		&newUser.Id,
		&newUser.FullName,
		&newUser.Username,
		&newUser.Email,
	)
	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

func (r *Repository) GetUserCache(id int) (string, error) {
	res, err := r.rdConn.JSONGet(r.ctx, getUserProfileCacheKey(id), userProfilePath).Result()
	if err != nil {
		return "", err
	}
	return res, nil
}

func (r *Repository) SetUserCache(user *types.User) (string, error) {
	res, err := r.rdConn.JSONSet(r.ctx, getUserProfileCacheKey(user.Id), userProfilePath, user).Result()
	if err != nil {
		return "", err
	}

	_, err = r.rdConn.Expire(r.ctx, getUserProfileCacheKey(user.Id), userProfileExpirationTime).Result()
	if err != nil {
		return "", err
	}

	return res, nil
}

func (r *Repository) GetUserRecentTweetsCache(id int) (string, error) {
	res, err := r.rdConn.JSONGet(r.ctx, getRecentTweetsCacheKey(id), userProfileRecentTWeetsPath).Result()
	if err != nil {
		return "", err
	}
	return res, nil
}

func (r *Repository) SetUserRecentTweetsCache(user *types.User, tweets []types.Tweet) (string, error) {
	res, err := r.rdConn.JSONSet(r.ctx, getRecentTweetsCacheKey(user.Id), userProfileRecentTWeetsPath, tweets).Result()
	if err != nil {
		return "", err
	}

	_, err = r.rdConn.Expire(r.ctx, getUserProfileCacheKey(user.Id), userProfileExpirationTime).Result()
	if err != nil {
		return "", err
	}

	return res, nil
}

func (r *Repository) DeleteUserRecentTweetsCache(id int) error {
	_, err := r.rdConn.JSONDel(r.ctx, getRecentTweetsCacheKey(id), userProfileRecentTWeetsPath).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetFeed(id int, page int) (*types.Feed, error) {
	limit := 10
	offset := (page - 1) * limit

	query := `
		SELECT 
			t.id,
			t.content,
			t.created_at,
			t.modified_at,
			t.like_count,
			CASE 
				WHEN tl.user_id IS NOT NULL THEN TRUE ELSE FALSE 
			END AS is_liked

			u.id,
			u.username,
			u.full_name,
			u.profile_image
		FROM 
			tweets t
		LEFT JOIN 
			follows f ON t.user_id = f.following_id
		LEFT JOIN 
			likes tl ON t.id = tl.tweet_id AND tl.user_id = @id
		INNER JOIN 
			users u ON u.id = t.user_id
		WHERE 
			f.follower_id = @id OR t.user_id = @id
		ORDER BY 
			t.created_at DESC, t.user_id DESC
		LIMIT 
			@limit
		OFFSET 
			@offset
	`
	args := pgx.NamedArgs{
		"id":     id,
		"limit":  limit,
		"offset": offset,
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
			&tweet.CreatedAt,
			&tweet.ModifiedAt,
			&tweet.LikeCount,
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
	if len(tweets) < limit {
		nextPageId = nil
	} else {
		nextPage := page + 1
		nextPageId = &nextPage
	}

	feed := types.Feed{
		Tweets:     tweets,
		NextPageId: nextPageId,
	}

	return &feed, nil
}

func (r *Repository) GetRecentTweets(userId int, page int) ([]types.Tweet, error) {
	limit := 10
	offset := (page - 1) * limit
	query := `
		SELECT 
			t.id,
			t.content,
			t.created_at,
			t.modified_at,
			t.like_count,
			FALSE,

			u.id,
			u.username,
			u.full_name,
			u.profile_image
		FROM 
			tweets t
		INNER JOIN 
			users u ON u.id = t.user_id
		WHERE 
			t.user_id = @userId
		ORDER BY 
			t.created_at DESC
		LIMIT 
			@limit
		OFFSET 
			@offset
	`
	args := pgx.NamedArgs{
		"userId": userId,
		"limit":  limit,
		"offset": offset,
	}
	rows, err := r.pgConn.Query(r.ctx, query, args)
	if err != nil {
		return nil, err
	}

	var lastTenTweets []types.Tweet
	for rows.Next() {
		var tweet types.Tweet
		var user types.User

		err := rows.Scan(
			&tweet.Id,
			&tweet.Content,
			&tweet.CreatedAt,
			&tweet.ModifiedAt,
			&tweet.LikeCount,
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
		lastTenTweets = append(lastTenTweets, tweet)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return lastTenTweets, nil
}

func (r *Repository) GetTweetsInteractions(userId int, tweetsId []int) ([]types.TweetInteraction, error) {
	query := `
		SELECT 
			tweet_id, 
			CASE 
				WHEN user_id = @userId THEN TRUE ELSE FALSE 
			END as is_liked
        FROM 
			likes
        WHERE 
			tweet_id = ANY(@tweetsId)
        	AND user_id = @userId
	`
	args := pgx.NamedArgs{
		"userId":   userId,
		"tweetsId": tweetsId,
	}
	rows, err := r.pgConn.Query(r.ctx, query, args)
	if err != nil {
		return nil, err
	}

	var tweetsInteraction []types.TweetInteraction
	for rows.Next() {
		var tweetInteraction types.TweetInteraction

		err := rows.Scan(
			&tweetInteraction.TweetId,
			&tweetInteraction.IsLiked,
		)
		if err != nil {
			return nil, err
		}

		tweetsInteraction = append(tweetsInteraction, tweetInteraction)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tweetsInteraction, nil
}

func (r *Repository) IsUserExistByEmail(email string) (bool, error) {
	var isUserExist bool
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE email=@email)`
	args := pgx.NamedArgs{
		"email": email,
	}

	err := r.pgConn.QueryRow(r.ctx, query, args).Scan(&isUserExist)
	if err != nil {
		return false, err
	}

	return isUserExist, nil
}

func (r *Repository) FindById(id int) (*types.User, error) {
	var user types.User
	query := `SELECT id, username, email, password, follower_count, following_count FROM users WHERE id=@id`
	args := pgx.NamedArgs{
		"id": id,
	}

	err := r.pgConn.QueryRow(r.ctx, query, args).Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.FollowerCount, &user.FollowingCount)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) FindByUsername(username string) (*types.User, error) {
	query := `
		SELECT 
			id,
			username,
			full_name,
			email,
			profile_image,
			follower_count,
			following_count
		FROM 
			users
		WHERE 
			username = @username
	`
	args := pgx.NamedArgs{
		"username": username,
	}

	var user types.User
	err := r.pgConn.QueryRow(r.ctx, query, args).Scan(
		&user.Id,
		&user.Username,
		&user.FullName,
		&user.Email,
		&user.ProfileImage,
		&user.FollowerCount,
		&user.FollowingCount,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) FindByEmail(email string) (*types.User, error) {
	var user types.User
	query := `
		SELECT 
			id, 
			username, 
			full_name, 
			email, 
			password, 
			profile_image,
			follower_count,
			following_count
		FROM 
			users 
		WHERE 
			email = @email
	`
	args := pgx.NamedArgs{
		"email": email,
	}

	err := r.pgConn.QueryRow(r.ctx, query, args).Scan(
		&user.Id,
		&user.Username,
		&user.FullName,
		&user.Email,
		&user.Password,
		&user.ProfileImage,
		&user.FollowerCount,
		&user.FollowingCount,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) IsFollowed(followerId int, followingId int) (bool, error) {
	var isFollowed bool
	query := `SELECT EXISTS (SELECT 1 FROM follows WHERE follower_id=@followerId AND following_id=@followingId)`
	args := pgx.NamedArgs{
		"followerId":  followerId,
		"followingId": followingId,
	}

	err := r.pgConn.QueryRow(r.ctx, query, args).Scan(&isFollowed)
	if err != nil {
		return false, err
	}

	return isFollowed, nil
}

func (r *Repository) FollowOtherUser(followerId int, followingId int) error {
	tx, err := r.pgConn.Begin(r.ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(r.ctx)

	query := `
		INSERT INTO follows (follower_id, following_id)
		VALUES (@follower_id, @following_id)
	`
	args := pgx.NamedArgs{
		"follower_id":  followerId,
		"following_id": followingId,
	}
	_, err = r.pgConn.Exec(r.ctx, query, args)
	if err != nil {
		return err
	}

	query = `
		UPDATE users
		SET follower_count = follower_count + 1
		WHERE id = @following_id
	`
	args = pgx.NamedArgs{
		"following_id": followingId,
	}
	_, err = r.pgConn.Exec(r.ctx, query, args)
	if err != nil {
		return err
	}

	query = `
		UPDATE users
		SET following_count = following_count + 1
		WHERE id = @follower_id
	`
	args = pgx.NamedArgs{
		"follower_id": followerId,
	}
	_, err = r.pgConn.Exec(r.ctx, query, args)
	if err != nil {
		return err
	}

	err = tx.Commit(r.ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UnfollowOtherUser(followerId int, followingId int) error {
	tx, err := r.pgConn.Begin(r.ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(r.ctx)

	query := `
		DELETE FROM follows
		WHERE follower_id = @follower_id
		AND following_id = @following_id
	`
	args := pgx.NamedArgs{
		"follower_id":  followerId,
		"following_id": followingId,
	}
	cmdTag, err := r.pgConn.Exec(r.ctx, query, args)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return nil
	}

	query = `
		UPDATE users
		SET follower_count = follower_count - 1
		WHERE id = @following_id
	`
	args = pgx.NamedArgs{
		"following_id": followingId,
	}
	_, err = r.pgConn.Exec(r.ctx, query, args)
	if err != nil {
		return err
	}

	query = `
		UPDATE users
		SET following_count = following_count - 1
		WHERE id = @follower_id
	`
	args = pgx.NamedArgs{
		"follower_id": followerId,
	}
	_, err = r.pgConn.Exec(r.ctx, query, args)
	if err != nil {
		return err
	}

	err = tx.Commit(r.ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) ClearRedisCache() error {
	err := r.rdConn.FlushAll(r.ctx).Err()
	if err != nil {
		return err
	}
	return nil
}
