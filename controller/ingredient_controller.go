package controller

import (
	"errors"
	"fmt"

	"github.com/denisyao1/welsh-academy-api/exception"
	"github.com/denisyao1/welsh-academy-api/model"
	"github.com/denisyao1/welsh-academy-api/service"
	"github.com/gofiber/fiber/v2"
)

// IngredientController contains methods to route ingredient related requests.
type IngredientController struct {
	BaseController
	service service.IngredientService
}

// NewIngredientController returns new IngredientController object.
func NewIngredientController(service service.IngredientService) IngredientController {
	return IngredientController{service: service}
}

//	CreateIngredient creates new ingredient.
//
// @Summary      Create ingredient
// @Description  Create an ingredient.
// @Description
// @Description  Require Admin Role.
// @Param request body schema.Ingredient true "Ingredient object"
// @Tags         Ingredients
// @Produce      json
// @Success      200 {object} model.Ingredient
// @Failure      400 {object} exception.ErrValidation
// @Failure      409 {object} ErrMessage
// @Failure      500
// @Security JWT
// @Router       /ingredients [post]
func (c IngredientController) CreateIngredient(ctx *fiber.Ctx) error {
	var ingredient model.Ingredient

	if err := ctx.BodyParser(&ingredient); err != nil {
		return ctx.Status(BadRequest).JSON(NewErrMessage("Failed to read request body."))
	}

	validationErr := c.service.Validate(ingredient)
	if validationErr.Field != "" {
		return ctx.Status(BadRequest).JSON(validationErr)
	}

	err := c.service.Create(&ingredient)
	if err != nil {
		if errors.Is(err, exception.ErrDuplicateKey) {
			message := fmt.Sprintf("An ingredient named '%s' already exists.", ingredient.Name)
			return ctx.Status(Conflict).JSON(NewErrMessage(message))
		}
		return c.HandleUnExpetedError(err, ctx)
	}

	return ctx.Status(Created).JSON(ingredient)
}

//	ListIngredients lists all ingredients.
//
// @Summary      List ingredients
// @Description  List ingredients.
// @Tags         Ingredients
// @Produce      json
// @Success      200 {array} model.Ingredient
// @Failure      500
// @Security JWT
// @Router       /ingredients [get]
func (c IngredientController) ListIngredients(ctx *fiber.Ctx) error {

	ingredients, err := c.service.FindAll()
	if err != nil {
		return c.HandleUnExpetedError(err, ctx)
	}
	return ctx.Status(OK).JSON(Map{"count": len(ingredients), "ingredients": ingredients})
}
