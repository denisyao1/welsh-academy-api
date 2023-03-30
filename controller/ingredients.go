package controller

import (
	"errors"
	"fmt"

	"github.com/denisyao1/welsh-academy-api/exception"
	"github.com/denisyao1/welsh-academy-api/model"
	"github.com/denisyao1/welsh-academy-api/service"
	"github.com/gofiber/fiber/v2"
)

type IngredientController struct {
	service service.IngredientService
}

func NewIngredientController(service service.IngredientService) IngredientController {
	return IngredientController{service: service}
}

func (c IngredientController) CreateIngredient(ctx *fiber.Ctx) error {
	var ingredient model.Ingredient

	if err := ctx.BodyParser(&ingredient); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Failed to read request body"})
	}
	validationErr := c.service.Validate(ingredient)
	if validationErr != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": validationErr})
	}

	err := c.service.Create(&ingredient)
	if err != nil {

		if errors.Is(err, exception.ErrDuplicateKey) {
			message := fmt.Sprintf("An ingredient named '%s' already exists.", ingredient.Name)
			return ctx.Status(409).JSON(fiber.Map{"error": message})
		}
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})

	}

	return ctx.Status(201).JSON(ingredient)
}

func (c IngredientController) ListAllIngredients(ctx *fiber.Ctx) error {

	ingredients, err := c.service.FindAll()

	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err})
	}

	return ctx.Status(200).JSON(ingredients)

}
