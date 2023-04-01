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
	BaseController
	service service.RecipeService
}

func NewRecipeController(service service.RecipeService) RecipeController {
	return RecipeController{service: service}
}

func (c RecipeController) CreateRecipe(ctx *fiber.Ctx) error {
	var recipe model.Recipe

	if err := ctx.BodyParser(&recipe); err != nil {
		return ctx.Status(BadRequest).JSON(Map{"error": "Failed to read request body"})
	}

	validationErrs := c.service.Validate(&recipe)
	if validationErrs != nil {
		if len(validationErrs) == 1 {
			return ctx.Status(BadRequest).JSON(Map{"error": validationErrs[0]})
		}
		return ctx.Status(BadRequest).JSON(Map{"errors": validationErrs})
	}

	err := c.service.Create(&recipe)
	if err != nil {
		if errors.Is(err, exception.ErrDuplicateKey) {
			message := fmt.Sprintf("A recipe named '%s' already exists.", recipe.Name)
			return ctx.Status(Conflict).JSON(Map{"error": message})
		}
		return c.HandleUnExpetedError(err, ctx)
	}

	return ctx.Status(Created).JSON(recipe)
}

func (c RecipeController) ListRecipes(ctx *fiber.Ctx) error {
	ingredientQuery := new(schema.IngredientQuerySchema)

	errQuery := ctx.QueryParser(ingredientQuery)
	if errQuery != nil {
		ctx.Status(BadRequest).JSON(Map{"error": "Failed to read request query string"})
	}

	ingredientNames := ingredientQuery.Ingredients
	recipes, err := c.service.ListAllPossible(ingredientNames)
	if err != nil {
		return c.HandleUnExpetedError(err, ctx)
	}

	return ctx.Status(OK).JSON(Map{"count": len(recipes), "recipes": recipes})
}

func (c RecipeController) FlagOrUnflag(ctx *fiber.Ctx) error {
	// get userID
	userID, err := c.GetConnectedUserID(ctx)
	if err != nil {
		return ctx.Status(BadRequest).JSON(Map{"error": exception.ErrMalFormedJWT.Error()})
	}

	// get recipe from path params ID
	recipeID, err := c.ConvertParamToInt("id", ctx)
	if err != nil {
		msg := "recipe not found"
		return ctx.Status(BadRequest).JSON(Map{"error": msg})
	}

	// add or remove recipe from user favorite
	message, err := c.service.AddRemoveFavorite(userID, recipeID)
	if err != nil {
		if err == exception.ErrRecordNotFound {
			return ctx.Status(BadRequest).JSON(Map{"error": "recipe " + err.Error()})
		}
		return c.HandleUnExpetedError(err, ctx)
	}
	return ctx.Status(OK).JSON(Map{"message": message})
}

func (c RecipeController) ListUserFavorites(ctx *fiber.Ctx) error {
	// get userID
	userID, err := c.GetConnectedUserID(ctx)
	if err != nil {
		return ctx.Status(BadRequest).JSON(Map{"error": exception.ErrMalFormedJWT.Error()})
	}

	// return user favorites recipes from
	recipes, err := c.service.FindUserFavorites(userID)
	if err != nil {
		return c.HandleUnExpetedError(err, ctx)
	}

	return ctx.Status(OK).JSON(Map{"count": len(recipes), "recipes": recipes})
}
