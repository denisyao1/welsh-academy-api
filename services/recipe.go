package services

import (
	"fmt"

	"github.com/denisyao1/welsh-academy-api/exceptions"
	"github.com/denisyao1/welsh-academy-api/models"
	"github.com/denisyao1/welsh-academy-api/repositories"
	"github.com/denisyao1/welsh-academy-api/utils"
)

type RecipeService interface {
	ValidateAndTransform(recipe *models.Recipe) []error
	Create(recipe *models.Recipe) error
	transform(recipe *models.Recipe) []error
}

type recipeService struct {
	recipeRepo     repositories.RecipeRepository
	ingredientRepo repositories.IngredientRepository
}

func NewRecipeService(recipeRepo repositories.RecipeRepository, ingredientRepo repositories.IngredientRepository) RecipeService {
	return &recipeService{recipeRepo: recipeRepo, ingredientRepo: ingredientRepo}
}

func (s recipeService) ValidateAndTransform(recipe *models.Recipe) []error {
	var newExeption = exceptions.NewValidationError

	// recipe name must be non null
	if recipe.Name == "" {
		return []error{newExeption("name", "the name is required")}
	}

	// recipe ingredients slice must contains at least one element
	if len(recipe.Ingredients) == 0 {
		return []error{newExeption("ingredients", "recipe must contains a least one ingredient")}

	}

	// recipe ingredients slice  must not contains duplicate
	noDuplicate := utils.SliceHasNoDuplicate(recipe.Ingredients)
	if !noDuplicate {
		return []error{newExeption("ingredients", "recipe ingredients contains duplicate")}
	}

	return s.transform(recipe)
}

func (s recipeService) transform(recipe *models.Recipe) []error {
	var names []string

	for _, i := range recipe.Ingredients {
		names = append(names, i.Name)
	}

	dbIngredients, err := s.ingredientRepo.FindNamed(names)
	if err != nil {
		return []error{err}
	}

	if len(names) == len(dbIngredients) {
		recipe.Ingredients = dbIngredients
		return nil
	}

	var errs []error
	var dbNames []string
	var newErr = exceptions.NewValidationError

	for _, elm := range dbIngredients {
		dbNames = append(dbNames, elm.Name)
	}

	for _, name := range names {
		if !utils.Contains(name, dbNames) {
			errs = append(errs, newErr("ingredients", fmt.Sprintf("'%s' is not a valid ingredient", name)))
		}
	}
	return errs
}

func (s recipeService) Create(recipe *models.Recipe) error {
	ok, err := s.recipeRepo.CheckIfNotCreated(*recipe)

	if err != nil {
		return err
	}

	if !ok {
		return exceptions.ErrDuplicateKey
	}
	return s.recipeRepo.Create(recipe)
}
