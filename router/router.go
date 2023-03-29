package router

import (
	"log"
	"os"

	"github.com/denisyao1/welsh-academy-api/controllers"
	"github.com/gofiber/fiber/v2"
)

type Router struct {
	ingredientController controllers.IngredientController
	recipeController     controllers.RecipeController
}

func New(
	ingredientController controllers.IngredientController,
	recipeController controllers.RecipeController,
) *Router {
	return &Router{
		ingredientController: ingredientController,
		recipeController:     recipeController,
	}
}

func (r Router) InitRoutes(app *fiber.App) {
	app.Get("/health", controllers.HealthCheck)
	app.Post("/ingredients", r.ingredientController.CreateIngredient)
	app.Get("/ingredients", r.ingredientController.ListAllIngredients)
	app.Post("/recipes", r.recipeController.CreateRecipe)
	app.Get("/recipes", r.recipeController.ListRecipes)
}

func (r Router) Start(app *fiber.App) {
	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
