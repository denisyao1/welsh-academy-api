package services

import (
	"os"
	"strconv"
	"time"

	"github.com/denisyao1/welsh-academy-api/exceptions"
	"github.com/denisyao1/welsh-academy-api/models"
	"github.com/denisyao1/welsh-academy-api/repositories"
	"github.com/denisyao1/welsh-academy-api/schemas"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	ValidateUserCreation(userSchema schemas.CreateUserSchema) []error
	validateCredentials(loginSchema schemas.LoginSchema) (models.User, error)
	CreateUser(userSchema schemas.CreateUserSchema) (models.User, error)
	CreateAccessToken(loginSchema schemas.LoginSchema) (string, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s userService) ValidateUserCreation(userSchema schemas.CreateUserSchema) []error {
	var newErrValidation = exceptions.NewValidationError
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

func (s userService) CreateUser(userSchema schemas.CreateUserSchema) (models.User, error) {

	user := models.User{Username: userSchema.Username, IsAdmin: userSchema.IsAdmin}

	// Check if username is already used in DB
	ok, checkErr := s.repo.CheckIfNotCreated(user)

	if checkErr != nil {
		return user, checkErr
	}

	if !ok {
		return user, exceptions.ErrDuplicateKey
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

func (s userService) validateCredentials(loginSchema schemas.LoginSchema) (models.User, error) {
	var user models.User

	if loginSchema.Username == "" || loginSchema.Pasword == "" {
		return user, exceptions.ErrInvalidCredentials
	}

	// username must be at least 3 characters long; password at least 4 characters long
	if len([]rune(loginSchema.Username)) < 3 || len([]rune(loginSchema.Username)) < 4 {
		return user, exceptions.ErrInvalidCredentials
	}

	user.Username = loginSchema.Username
	err := s.repo.GetUserByUsername(&user)
	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, exceptions.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginSchema.Pasword))

	if err != nil {
		return user, exceptions.ErrInvalidCredentials
	}

	return user, nil
}

func (s userService) CreateAccessToken(loginSchema schemas.LoginSchema) (string, error) {
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
