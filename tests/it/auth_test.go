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

func TestAuthRegister_Ok(t *testing.T) {
	ResetAndSeed()
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

func TestAuthRegister_EmailAlreadyExist(t *testing.T) {
	ResetAndSeed()
	userRepository := user.NewRepository(ctx, pgConn, rdConn)
	userService := user.NewService(ctx, cfg, userRepository)

	testUser := models.User{
		Email:    validUser.Email,
		Password: faker.Password(),
	}

	_, err := userService.CreateUser(testUser)
	assert.EqualError(t, err, errmsg.EMAIL_ALREADY_EXIST)
}

func TestAuthLogin_Ok(t *testing.T) {
	ResetAndSeed()
	userRepository := user.NewRepository(ctx, pgConn, rdConn)
	userService := user.NewService(ctx, cfg, userRepository)

	user, err := userService.CheckUserCredential(validUser.Email, validUser.Password)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, user.Email, validUser.Email)
	assert.Equal(t, user.Username, validUser.Username)
	assert.Equal(t, user.FullName, validUser.FullName)
}

func TestAuthLogin_UserNotFound(t *testing.T) {
	ResetAndSeed()
	userRepository := user.NewRepository(ctx, pgConn, rdConn)
	userService := user.NewService(ctx, cfg, userRepository)

	_, err := userService.CheckUserCredential(faker.Email(), faker.Password())
	assert.EqualError(t, err, errmsg.USER_NOT_FOUND)
}

func TestAuthLogin_WrongCredential(t *testing.T) {
	ResetAndSeed()
	userRepository := user.NewRepository(ctx, pgConn, rdConn)
	userService := user.NewService(ctx, cfg, userRepository)

	_, err := userService.CheckUserCredential(validUser.Email, faker.Password())
	assert.EqualError(t, err, errmsg.WRONG_CREDENTIAL)
}
