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
		// return &models.AppError{Error: err, Message: "Failed to create account", Code: 500}
		return nil, err
	}
	return &newUser, nil
}
