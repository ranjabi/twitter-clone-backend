package user

import (
	"context"
	"errors"
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
		return nil, errors.New("failed to create account")
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
		return false, errors.New("failed to check user account")
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
		return nil, errors.New("failed to get user credential")
	}

	return &user, nil
}
