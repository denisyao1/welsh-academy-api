package repository

import (
	"github.com/denisyao1/welsh-academy-api/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	CheckIfNotCreated(user model.User) (bool, error)
	Create(user *model.User) error
	GetUserByUsername(user *model.User) error
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

func (r userRepo) Create(user *model.User) error {

	err := r.db.Create(user).Error

	return err
}

func (r userRepo) CheckIfNotCreated(user model.User) (bool, error) {
	var userB model.User
	result := r.db.Where("username=?", user.Username).Find(&userB)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected == 0, nil
}

func (r userRepo) GetUserByUsername(user *model.User) error {

	return r.db.Find(user).Error
}
