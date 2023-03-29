package controllers

import (
	"errors"
	"fmt"

	"github.com/denisyao1/welsh-academy-api/exceptions"
	"github.com/denisyao1/welsh-academy-api/models"
	"github.com/denisyao1/welsh-academy-api/services"
	"github.com/gofiber/fiber/v2"
)

type RecipeController struct {
	service services.RecipeService
}

func NewRecipeController(service services.RecipeService) RecipeController {
	return RecipeController{service: service}
}

func (c RecipeController) CreateRecipe(ctx *fiber.Ctx) error {
	var recipe models.Recipe

	ctx.BodyParser(&recipe)
	validationErrs := c.service.ValidateAndTransform(&recipe)
	if validationErrs != nil {
		if len(validationErrs) == 1 {
			return ctx.Status(400).JSON(fiber.Map{"error": validationErrs[0]})
		}
		return ctx.Status(400).JSON(fiber.Map{"errors": validationErrs})
	}

	err := c.service.Create(&recipe)
	if err != nil {
		if errors.Is(err, exceptions.ErrDuplicateKey) {
			message := fmt.Sprintf("A recipe named '%s' already exists.", recipe.Name)
			return ctx.Status(409).JSON(fiber.Map{"error": message})
		}
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})

	}

	return ctx.Status(201).JSON(recipe)
}
