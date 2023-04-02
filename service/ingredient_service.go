package service

import (
	"github.com/denisyao1/welsh-academy-api/exception"
	"github.com/denisyao1/welsh-academy-api/model"
	"github.com/denisyao1/welsh-academy-api/repository"
)

type IngredientService interface {
	Validate(ingredient model.Ingredient) exception.ErrValidation
	Create(ingredient *model.Ingredient) error
	FindAll() ([]model.Ingredient, error)
}

func NewIngredientService(repository repository.IngredientRepository) IngredientService {
	return &ingredientService{repo: repository}
}

type ingredientService struct {
	repo repository.IngredientRepository
}

func (s ingredientService) Validate(ingredient model.Ingredient) exception.ErrValidation {
	if ingredient.Name == "" {
		return exception.NewValidationError("name", "the name is reqired")
	}
	return nil

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
