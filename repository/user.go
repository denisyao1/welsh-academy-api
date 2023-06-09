package repository

import (
	"errors"

	"github.com/denisyao1/welsh-academy-api/exception"
	"github.com/denisyao1/welsh-academy-api/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	// IsNotCreated returns true if user is not present in DB else false
	IsNotCreated(user model.User) (bool, error)

	// Create adds user to DB.
	Create(user *model.User) error

	// GetByUsername returns user model corresponding to username.
	GetByUsername(user *model.User) error

	// GetByID returns user from DB by its ID.
	GetByID(user *model.User) error

	// UpdatePassword updates user password.
	UpdatePassword(user *model.User) error
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

func (r userRepo) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r userRepo) IsNotCreated(user model.User) (bool, error) {
	var userB model.User
	err := r.db.Where("username=?", user.Username).First(&userB).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return true, nil
	}
	return false, err
}

func (r userRepo) GetByUsername(user *model.User) error {
	err := r.db.Where("username=?", user.Username).First(user).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return exception.ErrRecordNotFound
	}
	return err
}

func (r userRepo) GetByID(user *model.User) error {
	err := r.db.Where("id=?", user.ID).First(user).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return exception.ErrRecordNotFound
	}
	return err
}

func (r userRepo) UpdatePassword(user *model.User) error {
	err := r.db.Model(user).Update("password", user.Password).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return exception.ErrRecordNotFound
	}
	return err
}
