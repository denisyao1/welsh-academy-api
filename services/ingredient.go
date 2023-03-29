package services

import (
	"github.com/denisyao1/welsh-academy-api/exceptions"
	"github.com/denisyao1/welsh-academy-api/models"
	"github.com/denisyao1/welsh-academy-api/repositories"
)

type IngredientService interface {
	Validate(ingredient models.Ingredient) exceptions.ErrValidation
	Create(ingredient *models.Ingredient) error
	FindAll() ([]models.Ingredient, error)
}

func NewIngredientService(repository repositories.IngredientRepository) IngredientService {
	return &ingredientService{repo: repository}
}

type ingredientService struct {
	repo repositories.IngredientRepository
}

func (s ingredientService) Validate(ingredient models.Ingredient) exceptions.ErrValidation {
	if ingredient.Name == "" {
		return exceptions.NewValidationError("name", "the name is reqired")
	}
	return nil

}

func (s ingredientService) Create(ingredient *models.Ingredient) error {

	ok, err := s.repo.CheckIfNotCreated(ingredient)

	if err != nil {
		return err
	}

	if !ok {
		return exceptions.ErrDuplicateKey
	}

	return s.repo.Create(ingredient)

}

func (s ingredientService) FindAll() ([]models.Ingredient, error) {

	return s.repo.FindAll()
}
