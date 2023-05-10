package service

import (
	"github.com/denisyao1/welsh-academy-api/exception"
	"github.com/denisyao1/welsh-academy-api/model"
	"github.com/denisyao1/welsh-academy-api/repository"
)

// IngredientService contains business logic to save and retreive ingredients.
type IngredientService interface {
	// Validate validates user input
	Validate(ingredient model.Ingredient) exception.ErrValidation

	// Create adds new ingredient to database.
	//
	// It returns exception.ErrDuplicateKey if the recipe name is alredy used.
	Create(ingredient *model.Ingredient) error

	// FindAll returns all the ingredients from the database.
	FindAll() ([]model.Ingredient, error)
}

// NewIngredientService returns new IngredientService.
func NewIngredientService(repository repository.IngredientRepository) IngredientService {
	return &ingredientService{repo: repository}
}

type ingredientService struct {
	repo repository.IngredientRepository
}

func (s ingredientService) Validate(ingredient model.Ingredient) exception.ErrValidation {
	if ingredient.Name == "" {
		return exception.NewErrValidation("name", "the name is reqired")
	}
	return exception.ErrValidation{}

}

func (s ingredientService) Create(ingredient *model.Ingredient) error {
	// Check if ingredient is already in the database
	ok, err := s.repo.IsNotCreated(*ingredient)
	if err != nil {
		return err
	}

	if !ok {
		return exception.ErrDuplicateKey
	}

	return s.repo.Create(ingredient)

}

func (s ingredientService) FindAll() ([]model.Ingredient, error) {

	return s.repo.FindAll()
}
