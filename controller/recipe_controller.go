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

// RecipeCOntroller contains methods to route recipes related requests.
type RecipeController struct {
	BaseController
	service service.RecipeService
}

// NewRecipeController returns new recipe controller.
func NewRecipeController(service service.RecipeService) RecipeController {
	return RecipeController{service: service}
}

//	CreateRecipe creates new recipe.
//
// @Summary      Create recipe
// @Description  Create recipe.
// @Description
// @Description  Require Admin Role.
// @Param request body schema.Recipe true "Recipe object"
// @Tags         Recipes
// @Accept       json
// @Produce      json
// @Success      201 {object} model.Recipe
// @Failure      400 {array} exception.ErrValidation
// @Failure      401 {object} ErrMessage
// @Failure      409 {object} ErrMessage
// @Failure      500
// @Security JWT
// @Router       /recipes [post]
func (c RecipeController) CreateRecipe(ctx *fiber.Ctx) error {
	var recipe model.Recipe

	if err := ctx.BodyParser(&recipe); err != nil {
		return ctx.Status(BadRequest).JSON(NewErrMessage("Failed to read request body"))
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
			return ctx.Status(Conflict).JSON(NewErrMessage(message))
		}
		return c.HandleUnExpetedError(err, ctx)
	}

	return ctx.Status(Created).JSON(recipe)
}

//	ListRecipes Lists all possible recipes.
//
// @Summary      List all possible recipes
// @Description  List all possible recipes.
// @Param 		 ingredients   query  schema.IngredientQuery false "ingredients"
// @Tags         Recipes
// @Accept       json
// @Produce      json
// @Success      200 {object} schema.RecipesResponse
// @Failure      400 {object} ErrMessage
// @Failure      401 {object} ErrMessage
// @Failure      500
// @Security JWT
// @Router       /recipes [get]
func (c RecipeController) ListRecipes(ctx *fiber.Ctx) error {
	ingredientQuery := schema.IngredientQuery{}
	errQuery := ctx.QueryParser(&ingredientQuery)
	if errQuery != nil {
		ctx.Status(BadRequest).JSON(NewErrMessage("Failed to read request query string"))
	}
	ingredientNames := ingredientQuery.Ingredients
	recipes, err := c.service.ListAllPossible(ingredientNames)
	if err != nil {
		return c.HandleUnExpetedError(err, ctx)
	}

	return ctx.Status(OK).JSON(Map{"count": len(recipes), "recipes": recipes})
}

//	FlagOrUnflag add or remove a recipe to user favorites.
//
// @Summary      Flag or Unflag recipe
// @Description  Add or remove a recipe to your favorites.
// @Param 		 id   path  int true "recipe ID"
// @Tags         Recipes
// @Accept       json
// @Produce      json
// @Success      200 {object} Message
// @Failure      400 {object} ErrMessage
// @Failure      404 {object} ErrMessage
// @Failure      401 {object} ErrMessage
// @Failure      500
// @Security JWT
// @Router       /recipes/{id}/flag-unflag [post]
func (c RecipeController) FlagOrUnflag(ctx *fiber.Ctx) error {
	// get userID
	userID, err := c.GetConnectedUserID(ctx)
	if err != nil {
		return ctx.Status(BadRequest).JSON(NewErrMessage(exception.ErrMalFormedJWT.Error()))
	}

	// get recipe from path params ID
	recipeID, err := c.ConvertParamToInt("id", ctx)
	if err != nil {
		msg := "Failed to convert recipe id."
		return ctx.Status(BadRequest).JSON(NewErrMessage(msg))
	}

	// add or remove recipe from user favorite
	message, err := c.service.AddOrRemoveFavorite(userID, recipeID)
	if err != nil {
		if err == exception.ErrRecordNotFound {
			return ctx.Status(NotFound).JSON(Map{"error": "recipe " + err.Error()})
		}
		return c.HandleUnExpetedError(err, ctx)
	}
	return ctx.Status(OK).JSON(NewMessage(message))
}

//	ListUserFavorites list the connected user favorite recipes.
//
// @Summary      List favorite recipes
// @Description  list the connected user favorite recipes.
// @Tags         User Profile
// @Accept       json
// @Produce      json
// @Success      200 {object} schema.RecipesResponse
// @Failure      400 {object} ErrMessage
// @Failure      401 {object} ErrMessage
// @Failure      500
// @Security JWT
// @Router       /recipes/favorites [get]
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
