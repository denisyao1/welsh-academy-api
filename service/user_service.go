package service

import (
	"errors"
	"log"
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
	ValidateUserCreation(userSchema schema.User) []error

	// validateCredentials checks if login informations are valid
	validateCredentials(loginSchema schema.Login) (model.User, error)

	// Create create new user
	Create(userSchema schema.User) (model.User, error)

	// CreateAccessToken return new access token
	CreateAccessToken(loginSchema schema.Login) (string, error)

	// UpdatePaswword Updates connected user password.
	//
	// if it receives bad input, it can returns :
	//		- exception.ErrInvalidPassword
	//		- exception.ErrRecordNotFound
	//      - exception.ErrPasswordSame
	UpdatePaswword(userID int, newPwd schema.Password) error

	// GetInfos returns the connected user model object.
	//
	//it returns exception.ErrRecordNotFound if user is not found
	GetInfos(userID int) (model.User, error)

	// CreateDefaultAdmin adds default admin user to DB
	CreateDefaultAdmin()

	// CreateIfNotExist creates a user in the DB if it's not already created.
	CreateIfNotExist(user *model.User) error
}

type userService struct {
	repo       repository.UserRepository
	jwt_secret string
}

func NewUserService(repo repository.UserRepository, jwt_secret string) UserService {
	return &userService{repo: repo, jwt_secret: jwt_secret}
}

func (s userService) ValidateUserCreation(userSchema schema.User) []error {
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

func (s userService) Create(userSchema schema.User) (model.User, error) {

	user := model.User{Username: userSchema.Username, IsAdmin: userSchema.IsAdmin}

	// Check if username is already used in DB
	ok, checkErr := s.repo.IsNotCreated(user)

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

func (s userService) validateCredentials(loginSchema schema.Login) (model.User, error) {
	var user model.User

	if loginSchema.Username == "" || loginSchema.Password == "" {
		return user, exception.ErrInvalidCredentials
	}

	// username must be at least 3 characters long; password at least 4 characters long
	if len([]rune(loginSchema.Username)) < 3 || len([]rune(loginSchema.Password)) < 4 {
		return user, exception.ErrInvalidCredentials
	}

	user.Username = loginSchema.Username
	err := s.repo.GetByUsername(&user)
	if err != nil {
		if errors.Is(err, exception.ErrRecordNotFound) {
			return user, exception.ErrInvalidCredentials
		}
		return user, err
	}

	if user.ID == 0 {
		return user, exception.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginSchema.Password))

	if err != nil {
		// return user, exception.ErrInvalidCredentials
		return user, err
	}

	return user, nil
}

func (s userService) CreateAccessToken(loginSchema schema.Login) (string, error) {
	// validate user credentials
	user, err := s.validateCredentials(loginSchema)
	if err != nil {
		return "", err
	}

	role := model.RoleUser
	if user.IsAdmin {
		role = model.RoleAdmin
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"ID":   strconv.Itoa(int(user.ID)),
		"role": role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	encodedToken, err := token.SignedString([]byte(s.jwt_secret))
	if err != nil {
		return "", err
	}

	return encodedToken, nil
}

func (s userService) UpdatePaswword(userID int, newPwdSchema schema.Password) error {
	//check if password is valid
	if newPwdSchema.Password == "" || len([]rune(newPwdSchema.Password)) < 4 {
		return exception.ErrInvalidPassword
	}

	// retrieve user from database
	var user model.User
	user.ID = userID
	err := s.repo.GetByID(&user)
	if err != nil {
		return err
	}

	if user.ID == 0 {
		return exception.ErrRecordNotFound
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(newPwdSchema.Password))
	if err == nil {
		return exception.ErrPasswordSame
	}

	hash, hashErr := bcrypt.GenerateFromPassword([]byte(newPwdSchema.Password), 10)
	if hashErr != nil {
		return hashErr
	}

	user.Password = string(hash)

	return s.repo.UpdatePassword(&user)
}

func (s userService) GetInfos(userID int) (model.User, error) {
	var user model.User
	user.ID = userID
	err := s.repo.GetByID(&user)
	return user, err
}

func (s userService) CreateDefaultAdmin() {
	var admin model.User
	admin.Username = "admin"
	admin.Password = "admin"
	admin.IsAdmin = true

	// check if admin is already created
	// or created it if already exist
	err := s.CreateIfNotExist(&admin)
	if err != nil {
		// Unexpected error occured exist
		log.Fatalln("Unexpected error : ", err.Error())
		log.Fatal("Failed to create default admin user")
	}
	log.Println("Default admin user created succssefully.")
}

func (s userService) CreateIfNotExist(user *model.User) error {
	err := s.repo.GetByUsername(user)
	if err != nil && !errors.Is(err, exception.ErrRecordNotFound) {
		return err
	}
	if user.ID != 0 {
		return nil
	}
	hash, hashErr := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if hashErr != nil {
		return hashErr
	}
	user.Password = string(hash)
	return s.repo.Create(user)
}
