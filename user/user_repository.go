package user

import (
	"context"
	"fmt"
	"time"
	"twitter-clone-backend/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

const (
	userProfilePath             = "$"
	userProfileRecentTWeetsPath = "$.recentTweets"
	userProfileExpirationTime   = 10 * time.Minute
)

func getUserProfileCacheKey(id int) string {
	return fmt.Sprintf("user.id:%d", id)
}

type UserRepository struct {
	ctx    context.Context
	pgConn *pgxpool.Pool
	rdConn *redis.Client
}

func NewRepository(ctx context.Context, pgConn *pgxpool.Pool, rdConn *redis.Client) UserRepository {
	return UserRepository{ctx: ctx, pgConn: pgConn, rdConn: rdConn}
}

func (r *UserRepository) GetUserCache(id int) (string, error) {
	res, err := r.rdConn.JSONGet(r.ctx, getUserProfileCacheKey(id), userProfilePath).Result()
	if err != nil {
		return "", err
	}
	return res, nil
}

func (r *UserRepository) GetUserRecentTweetsCache(id int) (string, error) {
	res, err := r.rdConn.JSONGet(r.ctx, getUserProfileCacheKey(id), userProfileRecentTWeetsPath).Result()
	if err != nil {
		return "", err
	}
	return res, nil
}

func (r *UserRepository) DeleteUserRecentTweetsCache(id int) error {
	_, err := r.rdConn.JSONDel(r.ctx, getUserProfileCacheKey(id), userProfileRecentTWeetsPath).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) SetUserCache(user *models.User) (string, error) {
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
func (r *UserRepository) SetUserRecentTweetsCache(user *models.User, tweets []models.Tweet) (string, error) {
	res, err := r.rdConn.JSONSet(r.ctx, getUserProfileCacheKey(user.Id), userProfileRecentTWeetsPath, tweets).Result()
	if err != nil {
		return "", err
	}

	_, err = r.rdConn.Expire(r.ctx, getUserProfileCacheKey(user.Id), userProfileExpirationTime).Result()
	if err != nil {
		return "", err
	}

	return res, nil
}

func (r *UserRepository) GetFeed(id int, page int) ([]models.Tweet, error) {
	limit := 10
	offset := (page - 1) * limit

	query := `
		SELECT t.*, u.username as username,
			CASE WHEN tl.user_id IS NOT NULL THEN TRUE ELSE FALSE END AS is_liked
		FROM tweets t
		LEFT JOIN follows f ON t.user_id = f.following_id
		LEFT JOIN likes tl ON t.id = tl.tweet_id AND tl.user_id = @id
		INNER JOIN users u ON u.id = t.user_id
		WHERE f.follower_id = @id OR t.user_id = @id
		ORDER BY t.created_at DESC
		LIMIT @limit
		OFFSET @offset
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

	feed, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Tweet])
	if err != nil {
		return nil, err
	}

	return feed, nil
}

func (r *UserRepository) CreateUser(user models.User) (*models.User, error) {
	var newUser models.User
	query := `INSERT INTO users (username, email, password) VALUES (LOWER(@username), LOWER(@email), @password) RETURNING username, email`
	args := pgx.NamedArgs{
		"username": user.Username,
		"email":    user.Email,
		"password": string(user.Password),
	}

	err := r.pgConn.QueryRow(r.ctx, query, args).Scan(&newUser.Username, &newUser.Email)
	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

func (r *UserRepository) GetLastTenTweets(userId int) ([]models.Tweet, error) {
	query := `
		SELECT t.*, u.username, FALSE as is_liked
			FROM tweets t
			INNER JOIN users u ON u.id = t.user_id
			WHERE t.user_id = @userId
			ORDER BY t.created_at DESC
			LIMIT 10
	`
	args := pgx.NamedArgs{
		"userId": userId,
	}
	rows, err := r.pgConn.Query(r.ctx, query, args)
	if err != nil {
		return nil, err
	}

	lastTenTweets, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Tweet])
	if err != nil {
		return nil, err
	}

	return lastTenTweets, nil
}

func (r *UserRepository) GetLastTenTweetsInteractions(userId int, tweetsId []int) ([]models.TweetInteraction, error) {
	query := `
		SELECT tweet_id as tweet_id, 
			CASE WHEN user_id = @userId THEN TRUE ELSE FALSE END as is_liked
        FROM likes
        WHERE tweet_id = ANY(@tweetsId)
        AND user_id = @userId
		LIMIT 10
	`
	args := pgx.NamedArgs{
		"userId":   userId,
		"tweetsId": tweetsId,
	}
	rows, err := r.pgConn.Query(r.ctx, query, args)
	if err != nil {
		return nil, err
	}

	lastTenTweets, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.TweetInteraction])
	if err != nil {
		return nil, err
	}

	return lastTenTweets, nil
}

func (r *UserRepository) IsUserExistByEmail(email string) (bool, error) {
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

func (r *UserRepository) GetUserById(id int) (*models.User, error) {
	var user models.User
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

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, email, password, follower_count, following_count FROM users WHERE username=@username`
	args := pgx.NamedArgs{
		"username": username,
	}

	err := r.pgConn.QueryRow(r.ctx, query, args).Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.FollowerCount, &user.FollowingCount)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, email, password FROM users WHERE email=@email`
	args := pgx.NamedArgs{
		"email": email,
	}

	err := r.pgConn.QueryRow(r.ctx, query, args).Scan(&user.Id, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) IsFollowed(followerId int, followingId int) (bool, error) {
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

func (r *UserRepository) FollowOtherUser(followerId int, followingId int) error {
	query := `INSERT INTO follows (follower_id, following_id) VALUES (@follower_id, @following_id)`
	args := pgx.NamedArgs{
		"follower_id":  followerId,
		"following_id": followingId,
	}

	_, err := r.pgConn.Exec(r.ctx, query, args)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) UnfollowOtherUser(followerId int, followingId int) error {
	query := `DELETE FROM follows WHERE follower_id=@follower_id and following_id=@following_id`
	args := pgx.NamedArgs{
		"follower_id":  followerId,
		"following_id": followingId,
	}

	_, err := r.pgConn.Exec(r.ctx, query, args)
	if err != nil {
		return err
	}

	return nil
}
