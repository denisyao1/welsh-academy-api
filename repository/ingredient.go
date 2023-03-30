package repository

import (
	"github.com/denisyao1/welsh-academy-api/model"
	"gorm.io/gorm"
)

type IngredientRepository interface {
	Create(ingredient *model.Ingredient) error
	FindAll() ([]model.Ingredient, error)
	CheckIfNotCreated(ingredient model.Ingredient) (bool, error)
	FindNamed(names []string) ([]model.Ingredient, error)
}

type gormIngredientRepo struct {
	db *gorm.DB
}

func NewGormIngredientRepository(db *gorm.DB) IngredientRepository {
	return &gormIngredientRepo{db: db}
}

func (r gormIngredientRepo) Create(ingredient *model.Ingredient) error {
	result := r.db.Create(ingredient)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r gormIngredientRepo) FindAll() ([]model.Ingredient, error) {
	var ingredients []model.Ingredient
	results := r.db.Find(&ingredients)
	if results.Error != nil {
		return nil, results.Error
	}
	return ingredients, nil
}

func (r gormIngredientRepo) CheckIfNotCreated(ingredient model.Ingredient) (bool, error) {
	var ingredientB model.Ingredient
	result := r.db.Where("name=?", ingredient.Name).Find(&ingredientB)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected == 0, nil
}

func (r gormIngredientRepo) FindNamed(names []string) ([]model.Ingredient, error) {
	var ingredients []model.Ingredient

	result := r.db.Where("name IN ?", names).Find(&ingredients)
	if result.Error != nil {
		return nil, result.Error
	}

	return ingredients, nil
}
