package router

import (
	"log"
	"os"

	"github.com/denisyao1/welsh-academy-api/controller"
	"github.com/denisyao1/welsh-academy-api/middleware"
	"github.com/denisyao1/welsh-academy-api/model"
	"github.com/gofiber/fiber/v2"
)

type Router struct {
	ingredientController controller.IngredientController
	recipeController     controller.RecipeController
	userController       controller.UserController
	SigningKey           string
}

func New(
	ingredientController controller.IngredientController,
	recipeController controller.RecipeController,
	userController controller.UserController,
	signingKey string,

) *Router {
	return &Router{
		ingredientController: ingredientController,
		recipeController:     recipeController,
		userController:       userController,
		SigningKey:           signingKey,
	}
}

func (r Router) InitRoutes(app *fiber.App) {

	api := app.Group("/api/v1")

	// routes thant required no auth
	api.Get("/health", controller.HealthCheck)
	api.Post("/login", r.userController.Login)
	api.Get("/logout", r.userController.Logout)

	// required user auth routes
	api.Use(middleware.JwtWare(r.SigningKey, model.RoleUser))
	api.Get("/ingredients", r.ingredientController.ListIngredients)
	api.Get("/recipes", r.recipeController.ListRecipes)
	api.Post("/recipes/:id/flag-unflag", r.recipeController.FlagOrUnflag)
	api.Get("/recipes/favorites", r.recipeController.ListUserFavorites)
	api.Get("/users/my-infos", r.userController.GetInfos)
	api.Patch("/users/password-change", r.userController.UpdatePassword)

	// required admin auth routes
	api.Use(middleware.JwtWare(r.SigningKey, model.RoleAdmin))
	api.Post("/users", r.userController.Create)
	api.Post("/ingredients", r.ingredientController.CreateIngredient)
	api.Post("/recipes", r.recipeController.CreateRecipe)

}

func (r Router) Start(app *fiber.App) {
	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
