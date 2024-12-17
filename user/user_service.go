package user

import (
	"twitter-clone-backend/model"
	"twitter-clone-backend/models"
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

func (s Service) CreateUser(user model.User) (*model.User, *models.ServiceError) {
	isUserExist, err := s.repository.IsUserExistByEmail(user.Email)
	if err != nil {
		// message buat direturn ke response, err dari repository di lift up
		// &Error() has error from repository
		// &Message has business logic error message
		return nil, &models.ServiceError{Err: err, Message: "Failed to check user account"}

		// return nil, fmt.Errorf("failed to check user account", err)
	}
	if isUserExist {
		return nil, &models.ServiceError{Err: nil, Message: "Email is already used"}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return nil, &models.ServiceError{Err: err, Message: "Failed to hash password"}
	}

	user.Password = string(hashedPassword)
	newUser, err := s.repository.CreateUser(user)
	if err != nil {
		return nil, &models.ServiceError{Err: err, Message: "Failed to create account"}
	}

	return newUser, nil
}

func (s Service) CheckUserCredential(email string, password string) (*model.User, error) {
	isUserExist, err := s.repository.IsUserExistByEmail(email)
	if err != nil {
		return nil, &models.ServiceError{Err: err, Message: "Failed to check user account"}
	}
	if !isUserExist {
		// TODO: SEND ERROR CODE FROM THIS 401
		return nil, &models.ServiceError{Err: nil, Message: "User not found. Please create an account"}
	}

	user, err := s.repository.GetUserByEmail(email)
	if err != nil {
		return nil, &models.ServiceError{Err: err, Message: "Failed to get user credential"}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, &models.ServiceError{Err: err, Message: "Email/password is wrong"}
	}

	claims := jwt.MapClaims{
		"userId":   user.Id,
		"username": user.Username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(utils.JWT_SIGNATURE_KEY))
	if err != nil {
		return nil, &models.ServiceError{Err: err, Message: "Failed to sign token"}
	}

	user.Token = signedToken

	return user, nil
}
