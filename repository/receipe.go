package repository

import (
	"errors"

	"github.com/denisyao1/welsh-academy-api/exception"
	"github.com/denisyao1/welsh-academy-api/model"
	"gorm.io/gorm"
)

type RecipeRepository interface {
	Create(receipe *model.Recipe) error
	CheckIfNotCreated(recipe model.Recipe) (bool, error)
	FindAll() ([]model.Recipe, error)
	FindAllContainging(ingredientNames []string) ([]model.Recipe, error)
	GetByID(recipeID int) (model.Recipe, error)
	IsInUserFavorites(userID, recipeID int) (bool, error)
	DeleteFromFavorites(userID, recipeID int) error
	AddToFavorites(userID, recipeID int) error
	FindFavorites(userID int) ([]model.Recipe, error)
}

type gormRecipeRepo struct {
	db *gorm.DB
}

func NewGormRecipeRepository(db *gorm.DB) RecipeRepository {
	return &gormRecipeRepo{db: db}
}

func (r gormRecipeRepo) CheckIfNotCreated(recipe model.Recipe) (bool, error) {
	var recipeB model.Recipe
	err := r.db.Where("name=?", recipe.Name).First(&recipeB).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return true, nil
	}
	return false, err
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

func (r gormRecipeRepo) GetByID(recipeID int) (model.Recipe, error) {
	var recipe model.Recipe
	err := r.db.Where("id = ?", recipeID).First(&recipe).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return recipe, exception.ErrRecordNotFound
	}
	return recipe, err
}

func (r gormRecipeRepo) IsInUserFavorites(userID int, recipeID int) (bool, error) {
	var userFavorite model.UserFavorite
	err := r.db.Table("user_favorites").
		Where("user_id = ? and recipe_id = ?", userID, recipeID).
		First(&userFavorite).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return true, err
}

func (r gormRecipeRepo) DeleteFromFavorites(userID, recipeID int) error {
	var user_favorite model.UserFavorite
	user_favorite.UserID = userID
	user_favorite.RecipeID = recipeID

	return r.db.Table("user_favorites").
		Where("user_id = ? and recipe_id = ?", userID, recipeID).
		Delete(&user_favorite).Error

}

func (r gormRecipeRepo) AddToFavorites(userID, recipeID int) error {
	var userFavorite model.UserFavorite
	userFavorite.UserID = userID
	userFavorite.RecipeID = recipeID
	return r.db.Table("user_favorites").Create(&userFavorite).Error
}

func (r gormRecipeRepo) FindFavorites(userID int) ([]model.Recipe, error) {
	var recipes []model.Recipe

	subQuery := r.db.Table("recipes AS r").
		Select("r.id").
		Joins("INNER JOIN user_favorites f ON f.recipe_id=r.id").
		Joins("INNER JOIN users u ON u.id=f.user_id").
		Where("u.id = ?", userID)

	err := r.db.Model(&model.Recipe{}).
		Preload("Ingredients").
		Where("id in (?)", subQuery).
		Find(&recipes).Error

	return recipes, err
}
