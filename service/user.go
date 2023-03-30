package service

import (
	"os"
	"strconv"
	"time"

	"github.com/denisyao1/welsh-academy-api/exception"
	"github.com/denisyao1/welsh-academy-api/model"
	"github.com/denisyao1/welsh-academy-api/repository"
	"github.com/denisyao1/welsh-academy-api/schema"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	ValidateUserCreation(userSchema schema.CreateUserSchema) []error
	validateCredentials(loginSchema schema.LoginSchema) (model.User, error)
	CreateUser(userSchema schema.CreateUserSchema) (model.User, error)
	CreateAccessToken(loginSchema schema.LoginSchema) (string, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s userService) ValidateUserCreation(userSchema schema.CreateUserSchema) []error {
	var newErrValidation = exception.NewValidationError
	var errs []error

	if userSchema.Username == "" {
		errs = append(errs, newErrValidation("username", "username is required"))
	}

	// Username must be at least 3 characters long
	if len([]rune(userSchema.Username)) < 3 {
		errs = append(errs, newErrValidation("username", "username must be at least 3 characters long"))
	}

	if userSchema.Password == "" {
		errs = append(errs, newErrValidation("password", "password is required"))
	}

	// Password must be at least 4 characters long
	if len([]rune(userSchema.Password)) < 4 {
		errs = append(errs, newErrValidation("password", "password must be at least 4 characters long"))
	}

	return errs
}

func (s userService) CreateUser(userSchema schema.CreateUserSchema) (model.User, error) {

	user := model.User{Username: userSchema.Username, IsAdmin: userSchema.IsAdmin}

	// Check if username is already used in DB
	ok, checkErr := s.repo.CheckIfNotCreated(user)

	if checkErr != nil {
		return user, checkErr
	}

	if !ok {
		return user, exception.ErrDuplicateKey
	}

	// hash password
	hash, hashErr := bcrypt.GenerateFromPassword([]byte(userSchema.Password), 10)
	if hashErr != nil {
		return user, hashErr
	}

	user.Password = string(hash)

	err := s.repo.Create(&user)

	return user, err
}

func (s userService) validateCredentials(loginSchema schema.LoginSchema) (model.User, error) {
	var user model.User

	if loginSchema.Username == "" || loginSchema.Pasword == "" {
		return user, exception.ErrInvalidCredentials
	}

	// username must be at least 3 characters long; password at least 4 characters long
	if len([]rune(loginSchema.Username)) < 3 || len([]rune(loginSchema.Username)) < 4 {
		return user, exception.ErrInvalidCredentials
	}

	user.Username = loginSchema.Username
	err := s.repo.GetUserByUsername(&user)
	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, exception.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginSchema.Pasword))

	if err != nil {
		return user, exception.ErrInvalidCredentials
	}

	return user, nil
}

func (s userService) CreateAccessToken(loginSchema schema.LoginSchema) (string, error) {
	// validate user credentials
	user, err := s.validateCredentials(loginSchema)
	if err != nil {
		return "", err
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"ID":  strconv.Itoa(int(user.ID)),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	encodedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return encodedToken, nil
}
