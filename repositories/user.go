package repositories

import (
	"github.com/denisyao1/welsh-academy-api/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CheckIfNotCreated(user models.User) (bool, error)
	Create(user *models.User) error
	GetUserByUsername(user *models.User) error
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

func (r userRepo) Create(user *models.User) error {

	err := r.db.Create(user).Error

	return err
}

func (r userRepo) CheckIfNotCreated(user models.User) (bool, error) {
	var userB models.User
	result := r.db.Where("username=?", user.Username).Find(&userB)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected == 0, nil
}

func (r userRepo) GetUserByUsername(user *models.User) error {

	return r.db.Find(user).Error
}
