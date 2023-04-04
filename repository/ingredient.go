package repository

import (
	"errors"

	"github.com/denisyao1/welsh-academy-api/model"
	"gorm.io/gorm"
)

type IngredientRepository interface {
	// Create adds new ingredient to DB.
	Create(ingredient *model.Ingredient) error

	// FindAll returns all ingredients from DB.
	FindAll() ([]model.Ingredient, error)

	// IsNotCreated returns true if the ingredient is not present in DB, else false.
	IsNotCreated(ingredient model.Ingredient) (bool, error)

	// FindNamed returns all ingredients those names equal the names parameters.
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

func (r gormIngredientRepo) IsNotCreated(ingredient model.Ingredient) (bool, error) {
	var ingredientB model.Ingredient
	err := r.db.Where("name=?", ingredient.Name).First(&ingredientB).Error

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return true, nil
	}
	return false, err
}

func (r gormIngredientRepo) FindNamed(names []string) ([]model.Ingredient, error) {
	var ingredients []model.Ingredient

	err := r.db.Where("name IN ?", names).Find(&ingredients).Error
	if err != nil {
		return nil, err
	}

	return ingredients, nil
}
