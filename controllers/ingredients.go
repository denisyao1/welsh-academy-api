package controllers

import (
	"github.com/denisyao1/welsh-academy-api/exceptions"
	"github.com/denisyao1/welsh-academy-api/models"
	"github.com/denisyao1/welsh-academy-api/services"
	"github.com/gofiber/fiber/v2"
)

type IngredientController struct {
	service services.IngredientService
}

func NewIngredientController(service services.IngredientService) IngredientController {
	return IngredientController{service: service}
}

func (c *IngredientController) CreateIngredient(ctx *fiber.Ctx) error {
	var ingredient models.Ingredient

	ctx.BodyParser(&ingredient)
	validationErr := c.service.Validate(ingredient)
	if validationErr != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": validationErr})
	}

	err := c.service.Create(&ingredient)
	if err != nil {
		if error, ok := err.(exceptions.DuplicateKeyError); ok {
			return ctx.Status(400).JSON(fiber.Map{"error": error.Error()})
		}
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})

	}

	return ctx.Status(201).JSON(ingredient)
}

func (c IngredientController) ListAllIngredients(ctx *fiber.Ctx) error {
	var ingredients []models.Ingredient
	err := c.service.FindAll(&ingredients)

	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err})
	}

	return ctx.Status(200).JSON(ingredients)

}
