package repositories

import (
	"github.com/denisyao1/welsh-academy-api/models"
	"gorm.io/gorm"
)

type RecipeRepository interface {
	Create(receipe *models.Recipe) error
	CheckIfNotCreated(recipe models.Recipe) (bool, error)
	// FindAllContainging(ingredients []models.Ingredient) ([]models.Recipe, error)
}

type gormRecipeRepo struct {
	db *gorm.DB
}

func NewGormRecipeRepository(db *gorm.DB) RecipeRepository {
	return &gormRecipeRepo{db: db}
}

func (r gormRecipeRepo) CheckIfNotCreated(recipe models.Recipe) (bool, error) {
	var recipeB models.Recipe
	result := r.db.Where("name=?", recipe.Name).Find(&recipeB)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected == 0, nil
}

func (r gormRecipeRepo) Create(recipe *models.Recipe) error {
	result := r.db.Create(recipe) // utiliser le omit pour ne pas créer les ingredients

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r gormRecipeRepo) FindAllContainging(ingredients []models.Ingredient) ([]models.Recipe, error) {

	return nil, nil

}
