package it

import (
	"net/http"
	"strings"
	"twitter-clone-backend/errmsg"
	"twitter-clone-backend/models"

	"github.com/go-faker/faker/v4"
)

func (s *TestSuite) TestAuthRegister_Ok() {
	testUser := models.User{
		Email:    faker.Email(),
		Password: faker.Password(),
	}

	newUser, err := s.userService.CreateUser(testUser)

	s.NoError(err)
	s.NotNil(newUser)
	s.Equal(strings.ToLower(testUser.Email), newUser.Email)
}

func (s *TestSuite) TestAuthRegister_EmailAlreadyExist() {
	duplicateUser := models.User{
		Email:    s.validUser.Email,
		Password: faker.Password(),
	}

	_, err := s.userService.CreateUser(duplicateUser)

	s.EqualError(err, errmsg.EMAIL_ALREADY_EXIST)
	s.IsType(&models.AppError{}, err)
	s.Equal(http.StatusConflict, err.(*models.AppError).GetCode())
}

func (s *TestSuite) TestAuthLogin_Ok() {
	user, err := s.userService.CheckUserCredential(s.validUser.Email, "password")

	s.NoError(err)
	s.NotNil(user)
	s.Equal(user.Email, s.validUser.Email)
	s.Equal(user.Username, s.validUser.Username)
	s.Equal(user.FullName, s.validUser.FullName)
}

func (s *TestSuite) TestAuthLogin_UserNotFound() {
	_, err := s.userService.CheckUserCredential(faker.Email(), faker.Password())

	s.EqualError(err, errmsg.USER_NOT_FOUND)
	s.IsType(&models.AppError{}, err)
	s.Equal(http.StatusNotFound, err.(*models.AppError).GetCode())
}

func (s *TestSuite) TestAuthLogin_WrongCredential() {
	_, err := s.userService.CheckUserCredential(s.validUser.Email, faker.Password())

	s.EqualError(err, errmsg.WRONG_CREDENTIAL)
	s.IsType(&models.AppError{}, err)
	s.Equal(http.StatusUnauthorized, err.(*models.AppError).GetCode())
}
