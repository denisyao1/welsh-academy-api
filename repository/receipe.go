package repository

import (
	"github.com/denisyao1/welsh-academy-api/model"
	"gorm.io/gorm"
)

type RecipeRepository interface {
	Create(receipe *model.Recipe) error
	CheckIfNotCreated(recipe model.Recipe) (bool, error)
	FindAll() ([]model.Recipe, error)
	FindAllContainging(ingredientNames []string) ([]model.Recipe, error)
}

type gormRecipeRepo struct {
	db *gorm.DB
}

func NewGormRecipeRepository(db *gorm.DB) RecipeRepository {
	return &gormRecipeRepo{db: db}
}

func (r gormRecipeRepo) CheckIfNotCreated(recipe model.Recipe) (bool, error) {
	var recipeB model.Recipe
	result := r.db.Where("name=?", recipe.Name).Find(&recipeB)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected == 0, nil
}

func (r gormRecipeRepo) Create(recipe *model.Recipe) error {
	err := r.db.Create(recipe).Error
	return err
}

func (r gormRecipeRepo) FindAll() ([]model.Recipe, error) {
	var recipes []model.Recipe

	err := r.db.Model(&model.Recipe{}).
		Preload("Ingredients").
		Find(&recipes).Error

	return recipes, err
}

func (r gormRecipeRepo) FindAllContainging(ingredientNames []string) ([]model.Recipe, error) {
	var recipes []model.Recipe

	subQuery := r.db.Table("recipes AS r").
		Select("r.id").
		Joins("INNER JOIN recipe_ingredients ri ON ri.recipe_id=r.id").
		Joins("INNER JOIN ingredients ing ON ing.id=ri.ingredient_id").
		Where("ing.name in ?", ingredientNames)

	err := r.db.Model(&model.Recipe{}).
		Preload("Ingredients").
		Where("id in (?)", subQuery).
		Find(&recipes).Error

	return recipes, err

}
