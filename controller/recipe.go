package controller

import (
	"errors"
	"fmt"

	"github.com/denisyao1/welsh-academy-api/exception"
	"github.com/denisyao1/welsh-academy-api/model"
	"github.com/denisyao1/welsh-academy-api/schema"
	"github.com/denisyao1/welsh-academy-api/service"
	"github.com/gofiber/fiber/v2"
)

type RecipeController struct {
	service service.RecipeService
}

func NewRecipeController(service service.RecipeService) RecipeController {
	return RecipeController{service: service}
}

func (c RecipeController) CreateRecipe(ctx *fiber.Ctx) error {
	var recipe model.Recipe

	if err := ctx.BodyParser(&recipe); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Failed to read request body"})
	}

	validationErrs := c.service.Validate(&recipe)
	if validationErrs != nil {
		if len(validationErrs) == 1 {
			return ctx.Status(400).JSON(fiber.Map{"error": validationErrs[0]})
		}
		return ctx.Status(400).JSON(fiber.Map{"errors": validationErrs})
	}

	err := c.service.Create(&recipe)
	if err != nil {
		if errors.Is(err, exception.ErrDuplicateKey) {
			message := fmt.Sprintf("A recipe named '%s' already exists.", recipe.Name)
			return ctx.Status(409).JSON(fiber.Map{"error": message})
		}
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})

	}

	return ctx.Status(201).JSON(recipe)
}

func (c RecipeController) ListRecipes(ctx *fiber.Ctx) error {
	ingredientQuery := new(schema.IngredientQuerySchema)
	errQuery := ctx.QueryParser(ingredientQuery)
	if errQuery != nil {
		ctx.Status(400).JSON(fiber.Map{"error": "Failed to read request query string"})
	}
	ingredientNames := ingredientQuery.Ingredients
	recipes, err := c.service.ListAllPossible(ingredientNames)
	if err != nil {
		ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(200).JSON(recipes)
}