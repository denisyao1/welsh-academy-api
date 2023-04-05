package repository

import (
	"testing"

	"github.com/denisyao1/welsh-academy-api/database"
	"github.com/denisyao1/welsh-academy-api/model"
	"github.com/stretchr/testify/assert"
)

func TestGetOrCreateIngredient(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	db, _ := database.NewInMemoryDB(false)
	db.Migrate(model.Ingredient{})

	ingredientRepo := NewGormIngredientRepository(db.GetDB())

	ingredient1, _ := ingredientRepo.GetOrCreate("ingredient1")
	ingredient1B, _ := ingredientRepo.GetOrCreate("ingredient1")

	assert.True(ingredient1.ID != 0, "ingredient ID should not be 0")
	assert.True(ingredient1.ID == ingredient1B.ID, "ingredient IDs should be the same")
	assert.True(ingredient1B.Name == ingredient1.Name, "ingredient Names should  be the same")
}

func TestGetOrCreateRecipe(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	db, _ := database.NewInMemoryDB(false)
	db.Migrate(model.Ingredient{}, model.Recipe{})

	ingredientRepo := NewGormIngredientRepository(db.GetDB())
	recipeRepo := NewGormRecipeRepository(db.GetDB())

	ingredient1, _ := ingredientRepo.GetOrCreate("ingredient1")
	ingredient2, _ := ingredientRepo.GetOrCreate("ingredient2")

	recipe1 := model.Recipe{
		Name:        "recipe1",
		Making:      "dummy",
		Ingredients: []model.Ingredient{ingredient1, ingredient2}}
	recipe2 := model.Recipe{
		Name:        "recipe1",
		Making:      "dummy",
		Ingredients: []model.Ingredient{ingredient1, ingredient2}}

	recipeRepo.GetOrCreate(&recipe1)
	recipeRepo.GetOrCreate(&recipe2)
	assert.True(recipe1.ID != 0, "recipe ID should not be 0")
	assert.True(recipe1.ID == recipe2.ID, "recipes ID should be the same")
	assert.True(recipe2.Name == recipe1.Name, "recipes Name should be the same ")
	ok := len(recipe1.Ingredients) == len(recipe2.Ingredients)
	assert.True(ok, "recipes ingredients length should be the same")
}
