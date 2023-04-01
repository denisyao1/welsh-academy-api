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
	BaseController
	service service.IngredientService
}

func NewIngredientController(service service.IngredientService) IngredientController {
	return IngredientController{service: service}
}

func (c IngredientController) CreateIngredient(ctx *fiber.Ctx) error {
	var ingredient model.Ingredient

	if err := ctx.BodyParser(&ingredient); err != nil {
		return ctx.Status(BadRequest).JSON(Map{"error": "Failed to read request body"})
	}
	validationErr := c.service.Validate(ingredient)
	if validationErr != nil {
		return ctx.Status(BadRequest).JSON(Map{"error": validationErr})
	}

	err := c.service.Create(&ingredient)
	if err != nil {

		if errors.Is(err, exception.ErrDuplicateKey) {
			message := fmt.Sprintf("An ingredient named '%s' already exists.", ingredient.Name)
			return ctx.Status(Conflict).JSON(Map{"error": message})
		}
		return c.HandleUnExpetedError(err, ctx)

	}

	return ctx.Status(Created).JSON(ingredient)
}

func (c IngredientController) ListAllIngredients(ctx *fiber.Ctx) error {

	ingredients, err := c.service.FindAll()

	if err != nil {
		c.HandleUnExpetedError(err, ctx)
	}

	return ctx.Status(OK).JSON(Map{"count": len(ingredients), "ingredients": ingredients})

}
