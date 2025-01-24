package it

import (
	"strings"
	"twitter-clone-backend/errmsg"
	"twitter-clone-backend/models"
	"twitter-clone-backend/usecases/user"

	"testing"

	"github.com/go-faker/faker/v4"
	_ "github.com/jackc/pgx/v5/stdlib" // for pgx sql driver

	"github.com/stretchr/testify/assert"
)

func TestRegister_Ok(t *testing.T) {
	userRepository := user.NewRepository(ctx, pgConn, rdConn)
	userService := user.NewService(ctx, cfg, userRepository)

	testUser := models.User{
		Email:    faker.Email(),
		Password: faker.Password(),
	}

	newUser, err := userService.CreateUser(testUser)
	assert.NoError(t, err)
	assert.NotNil(t, newUser)
	assert.Equal(t, strings.ToLower(testUser.Email), newUser.Email)
}

func TestRegister_EmailAlreadyExist(t *testing.T) {
	userRepository := user.NewRepository(ctx, pgConn, rdConn)
	userService := user.NewService(ctx, cfg, userRepository)

	testUser := models.User{
		Email:    validUser.Email,
		Password: faker.Password(),
	}

	_, err := userService.CreateUser(testUser)
	assert.EqualError(t, err, errmsg.EMAIL_ALREADY_EXIST)
}

func TestLogin_Ok(t *testing.T) {
	userRepository := user.NewRepository(ctx, pgConn, rdConn)
	userService := user.NewService(ctx, cfg, userRepository)

	user, err := userService.CheckUserCredential(validUser.Email, validUser.Password)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, user.Email, validUser.Email)
	assert.Equal(t, user.Username, validUser.Username)
	assert.Equal(t, user.FullName, validUser.FullName)
}

func TestLogin_UserNotFound(t *testing.T) {
	userRepository := user.NewRepository(ctx, pgConn, rdConn)
	userService := user.NewService(ctx, cfg, userRepository)

	_, err := userService.CheckUserCredential(faker.Email(), faker.Password())
	assert.EqualError(t, err, errmsg.USER_NOT_FOUND)
}

func TestLogin_WrongCredential(t *testing.T) {
	userRepository := user.NewRepository(ctx, pgConn, rdConn)
	userService := user.NewService(ctx, cfg, userRepository)

	_, err := userService.CheckUserCredential(validUser.Email, faker.Password())
	assert.EqualError(t, err, errmsg.WRONG_CREDENTIAL)
}
