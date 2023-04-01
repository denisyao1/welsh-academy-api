package service

import (
	"fmt"

	"github.com/denisyao1/welsh-academy-api/exception"
	"github.com/denisyao1/welsh-academy-api/model"
	"github.com/denisyao1/welsh-academy-api/repository"
	"github.com/denisyao1/welsh-academy-api/util"
)

type RecipeService interface {
	Validate(recipe *model.Recipe) []error
	Create(recipe *model.Recipe) error
	transform(recipe *model.Recipe) []error
	ListAllPossible(ingredientNames []string) ([]model.Recipe, error)
	AddRemoveFavorite(userID int, recipeID int) (string, error)
	FindUserFavorites(userID int) ([]model.Recipe, error)
}

type recipeService struct {
	recipeRepo     repository.RecipeRepository
	ingredientRepo repository.IngredientRepository
}

func NewRecipeService(recipeRepo repository.RecipeRepository, ingredientRepo repository.IngredientRepository) RecipeService {
	return &recipeService{recipeRepo: recipeRepo, ingredientRepo: ingredientRepo}
}

func (s recipeService) Validate(recipe *model.Recipe) []error {
	var newErrValidation = exception.NewValidationError
	var errs []error

	// recipe name must be non null
	if recipe.Name == "" {
		errs = append(errs, newErrValidation("name", "the name is required"))
	}

	// recipe making must not be empty

	if recipe.Making == "" {
		errs = append(errs, newErrValidation("making", "the making is required"))
	}

	// recipe ingredients slice must contains at least one element
	if len(recipe.Ingredients) == 0 {
		errs = append(errs, newErrValidation("ingredients", "recipe must contains a least one ingredient"))
	}

	// recipe ingredients slice  must not contains duplicate
	noDuplicate := util.SliceHasNoDuplicate(recipe.Ingredients)
	if !noDuplicate {
		errs = append(errs, newErrValidation("ingredients", "recipe ingredients contains duplicate"))
	}

	if len(errs) != 0 {
		return errs
	}

	return s.transform(recipe)
}

func (s recipeService) transform(recipe *model.Recipe) []error {
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
	var newErr = exception.NewValidationError

	for _, elm := range dbIngredients {
		dbNames = append(dbNames, elm.Name)
	}

	for _, name := range names {
		if !util.Contains(name, dbNames) {
			errs = append(errs, newErr("ingredients", fmt.Sprintf("'%s' is not a valid ingredient", name)))
		}
	}
	return errs
}

func (s recipeService) Create(recipe *model.Recipe) error {
	ok, err := s.recipeRepo.CheckIfNotCreated(*recipe)

	if err != nil {
		return err
	}

	if !ok {
		return exception.ErrDuplicateKey
	}
	return s.recipeRepo.Create(recipe)
}

func (s recipeService) ListAllPossible(ingredientNames []string) ([]model.Recipe, error) {
	if len(ingredientNames) == 0 {
		return s.recipeRepo.FindAll()
	}

	return s.recipeRepo.FindAllContainging(ingredientNames)
}

func (s recipeService) AddRemoveFavorite(userID int, recipeID int) (string, error) {
	// check if recipe existe in the DB
	_, err := s.recipeRepo.GetByID(recipeID)
	if err != nil {
		return "", err
	}

	// Check if recipe is already in user favorites
	exist, err := s.recipeRepo.IsInUserFavorites(userID, recipeID)
	if err != nil {
		return "", err
	}

	message := "recipe added to favorites"
	if exist {
		// remove recipe from user favorites
		err = s.recipeRepo.DeleteFromFavorites(userID, recipeID)
		message = "recipe removed from favorites"
	} else {
		err = s.recipeRepo.AddToFavorites(userID, recipeID)
	}

	return message, err
}

func (s recipeService) FindUserFavorites(userID int) ([]model.Recipe, error) {
	return s.recipeRepo.FindFavorites(userID)
}
