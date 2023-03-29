package repositories

import (
	"github.com/denisyao1/welsh-academy-api/models"
	"gorm.io/gorm"
)

type IngredientRepository interface {
	Create(ingredient *models.Ingredient) error
	FindAll() ([]models.Ingredient, error)
	CheckIfNotCreated(ingredient *models.Ingredient) (bool, error)
	FindNamed(names []string) ([]models.Ingredient, error)
}

type gormIngredientRepo struct {
	db *gorm.DB
}

func NewGormIngredientRepository(db *gorm.DB) IngredientRepository {
	return &gormIngredientRepo{db: db}
}

func (r gormIngredientRepo) Create(ingredient *models.Ingredient) error {
	result := r.db.Create(ingredient)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r gormIngredientRepo) FindAll() ([]models.Ingredient, error) {
	var ingredients []models.Ingredient
	results := r.db.Find(&ingredients)
	if results.Error != nil {
		return nil, results.Error
	}
	return ingredients, nil
}

func (r gormIngredientRepo) CheckIfNotCreated(ingredient *models.Ingredient) (bool, error) {
	var ingredientB models.Ingredient
	result := r.db.Where("name=?", ingredient.Name).Find(&ingredientB)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected == 0, nil
}

func (r gormIngredientRepo) FindNamed(names []string) ([]models.Ingredient, error) {
	var ingredients []models.Ingredient

	result := r.db.Where("name IN ?", names).Find(&ingredients)
	if result.Error != nil {
		return nil, result.Error
	}

	return ingredients, nil
}
