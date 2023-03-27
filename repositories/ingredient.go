package repositories

import (
	"github.com/denisyao1/welsh-academy-api/exceptions"
	"github.com/denisyao1/welsh-academy-api/models"
	"gorm.io/gorm"
)

type IngredientRepository interface {
	Create(ingredient *models.Ingredient) error
	FindAll(ingredients *[]models.Ingredient) error
}

type gormRepo struct {
	db *gorm.DB
}

func NewGormIngredientRepository(db *gorm.DB) IngredientRepository {
	return &gormRepo{db: db}
}

func (r gormRepo) Create(ingredient *models.Ingredient) error {
	result := r.db.Create(ingredient)

	if result.Error != nil {
		return exceptions.NewDuplicateKeyError(ingredient.Name)
	}
	return nil
}

func (r gormRepo) FindAll(ingredients *[]models.Ingredient) error {
	results := r.db.Find(ingredients)
	if results.Error != nil {
		return results.Error
	}
	return nil
}
