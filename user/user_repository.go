package user

import (
	"context"
	"fmt"
	"twitter-clone-backend/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	ctx    context.Context
	pgConn *pgxpool.Pool
	rdConn *redis.Client
}

func NewRepository(ctx context.Context, pgConn *pgxpool.Pool, rdConn *redis.Client) Repository {
	return Repository{ctx: ctx, pgConn: pgConn, rdConn: rdConn}
}

func (r Repository) GetUserProfileCache(id int) (string, error) {
	res, err := r.rdConn.JSONGet(r.ctx, fmt.Sprintf("userProfile-id:%d", id), "$").Result()
	if err != nil {
		return "", err
	}
	return res, nil
}
func (r Repository) SetUserProfileCache(user *models.User) (string, error) {
	res, err := r.rdConn.JSONSet(r.ctx, fmt.Sprintf("userProfile-id:%d", user.Id), "$", user).Result()
	if err != nil {
		return "", err
	}
	return res, nil
}

func (r Repository) CreateUser(user models.User) (*models.User, error) {
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

func (r Repository) GetLastTenTweets(userId int) ([]models.Tweet, error) {
	query := `
		SELECT * 
			FROM tweets
			ORDER BY created_at DESC
			LIMIT 10
	`
	rows, err := r.pgConn.Query(r.ctx, query)
	if err != nil {
		return nil, err
	}

	lastTenTweets, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Tweet])
	if err != nil {
		return nil, err
	}
	return lastTenTweets, nil
}

func (r Repository) IsUserExistByEmail(email string) (bool, error) {
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

func (r Repository) GetUserById(id int) (*models.User, error) {
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

func (r Repository) GetUserByEmail(email string) (*models.User, error) {
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

func (r Repository) FollowOtherUser(followerId int, followingId int) error {
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

func (r Repository) UnfollowOtherUser(followerId int, followingId int) error {
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
