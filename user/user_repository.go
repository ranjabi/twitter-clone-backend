package user

import (
	"context"
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

func (r Repository) CreateUser(user model.User) (*model.User, error) {
	var newUser model.User
	query := `INSERT INTO users (username, email, password) VALUES (LOWER(@username), LOWER(@email), @password) RETURNING username, email`
	args := pgx.NamedArgs{
		"username": user.Username,
		"email":    user.Email,
		"password": string(user.Password),
	}

	err := r.conn.QueryRow(r.ctx, query, args).Scan(&newUser.Username, &newUser.Email)
	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

func (r Repository) IsUserExistByEmail(email string) (bool, error) {
	var isUserExist bool
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE email=@email)`
	args := pgx.NamedArgs{
		"email": email,
	}

	err := r.conn.QueryRow(r.ctx, query, args).Scan(&isUserExist)
	if err != nil {
		return false, err
	}

	return isUserExist, nil
}

func (r Repository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	query := `SELECT id, username, email, password FROM users WHERE email=@email`
	args := pgx.NamedArgs{
		"email": email,
	}

	err := r.conn.QueryRow(r.ctx, query, args).Scan(&user.Id, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r Repository) FollowOtherUser(followersId int, followingId int) error {
	query := `INSERT INTO follows (followers_id, following_id) VALUES (@followers_id, @following_id)`
	args := pgx.NamedArgs{
		"followers_id": followersId,
		"following_id": followingId,
	}

	_, err := r.conn.Exec(r.ctx, query, args)
	if err != nil {
		return err
	}

	return nil
}
