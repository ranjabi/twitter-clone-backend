package it

// import (
// 	"net/http"
// 	"strings"
// 	"twitter-clone-backend/app"
// 	"twitter-clone-backend/errmsg"
// 	"twitter-clone-backend/types"

// 	"github.com/go-faker/faker/v4"
// )

// func (s *TestSuite) TestAuthRegister_Ok() {
// 	testUser := types.User{
// 		Email:    faker.Email(),
// 		Password: faker.Password(),
// 	}

// 	newUser, err := s.authService.Register(testUser)

// 	s.NoError(err)
// 	s.NotNil(newUser)
// 	s.Equal(strings.ToLower(testUser.Email), newUser.Email)
// }

// func (s *TestSuite) TestAuthRegister_EmailAlreadyExist() {
// 	duplicateUser := types.User{
// 		Email:    s.validUser.Email,
// 		Password: faker.Password(),
// 	}

// 	_, err := s.authService.Register(duplicateUser)

// 	s.ErrorContains(err, errmsg.EMAIL_ALREADY_EXIST)
// 	s.IsType(&app.Error{}, err)
// 	s.Equal(http.StatusConflict, err.(*app.Error).GetCode())
// }

// func (s *TestSuite) TestAuthLogin_Ok() {
// 	user, err := s.authService.Login(s.validUser.Email, "password")

// 	s.NoError(err)
// 	s.NotNil(user)
// 	s.Equal(user.Email, s.validUser.Email)
// 	s.Equal(user.Username, s.validUser.Username)
// 	s.Equal(user.FullName, s.validUser.FullName)
// }

// func (s *TestSuite) TestAuthLogin_UserNotFound() {
// 	_, err := s.authService.Login(faker.Email(), faker.Password())

// 	s.ErrorContains(err, errmsg.USER_NOT_FOUND)
// 	s.IsType(&app.Error{}, err)
// 	s.Equal(http.StatusNotFound, err.(*app.Error).GetCode())
// }

// func (s *TestSuite) TestAuthLogin_WrongCredential() {
// 	_, err := s.authService.Login(s.validUser.Email, faker.Password())

// 	s.ErrorContains(err, errmsg.WRONG_CREDENTIAL)
// 	s.IsType(&app.Error{}, err)
// 	s.Equal(http.StatusUnauthorized, err.(*app.Error).GetCode())
// }
