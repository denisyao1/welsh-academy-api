package services

import (
	"github.com/denisyao1/welsh-academy-api/exceptions"
	"github.com/denisyao1/welsh-academy-api/models"
	"github.com/denisyao1/welsh-academy-api/repositories"
	"github.com/go-playground/validator/v10"
)

type IngredientService interface {
	Validate(ingredient models.Ingredient) exceptions.ValidationError
	Create(ingredient *models.Ingredient) error
	FindAll(ingredients *[]models.Ingredient) error
}

func NewIngredientService(repository repositories.IngredientRepository) IngredientService {
	return &ingredientService{repo: repository, valide: *validator.New()}
}

type ingredientService struct {
	repo   repositories.IngredientRepository
	valide validator.Validate
}

func (s ingredientService) Validate(ingredient models.Ingredient) exceptions.ValidationError {
	err := s.valide.Struct(ingredient)
	if err != nil {
		detail := err.(validator.ValidationErrors)[0]
		var error = exceptions.NewValidationError(detail.StructField(), detail.Tag(), detail.Param())
		return error
	}

	return nil

}

func (s ingredientService) Create(ingredient *models.Ingredient) error {

	return s.repo.Create(ingredient)

}

func (s ingredientService) FindAll(ingredients *[]models.Ingredient) error {

	return s.repo.FindAll(ingredients)
}
