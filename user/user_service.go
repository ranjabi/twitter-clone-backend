package user

import (
    "errors"
    "twitter-clone-backend/model"
    "twitter-clone-backend/utils"

    jwt "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return Service{repository: repository}
}

func (s Service) CreateUser(user model.User) (*model.User, error) {
	isUserExist, err := s.repository.IsUserExistByEmail(user.Email)
	if err != nil {
		return nil, err
	}

	if isUserExist {
		return nil, errors.New("email is already used")
	}

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
    if err != nil {
        return nil, errors.New("failed to hash password")
    }

    user.Password = string(hashedPassword)
    newUser, err := s.repository.CreateUser(user)
    if err != nil {
        return nil, err
    }

    return newUser, nil
}

func (s Service) CheckUserCredential(email string, password string) (*model.User, error) {
    user, err := s.repository.GetUserByEmail(email)
    if err != nil {
        return nil, errors.New("failed to get user credential")
    } 
    if user == nil {
        // todo: how to set code 401 from here?
        return nil, errors.New("user not found. Please create an account")
    }

    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
        return nil, errors.New("email/password is wrong")
    }

    claims := jwt.MapClaims{
        "userId":   user.Id,
        "username": user.Username,
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signedToken, err := token.SignedString([]byte(utils.JWT_SIGNATURE_KEY))
    if err != nil {
        return nil, errors.New("failed to sign token")
    }

    user.Token = signedToken

    return user, nil
}

