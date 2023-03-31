package router

import (
	"log"
	"os"

	"github.com/denisyao1/welsh-academy-api/controller"
	"github.com/denisyao1/welsh-academy-api/middleware"
	"github.com/denisyao1/welsh-academy-api/model"

	// "github.com/gofiber/fiber/middleware"
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

	app.Get("/health", controller.HealthCheck)
	app.Post("/login", r.userController.Login)
	app.Get("/logout", r.userController.Logout)

	app.Use(middleware.JwtWare(r.SigningKey, model.RoleUser))
	app.Get("/ingredients", r.ingredientController.ListAllIngredients)
	app.Get("/recipes", r.recipeController.ListRecipes)
	app.Get("/infos", r.userController.GetInfos)
	app.Patch("/password-change", r.userController.UpdatePassword)

	app.Use(middleware.JwtWare(r.SigningKey, model.RoleAdmin))
	app.Post("/users", r.userController.Create)
	app.Post("/ingredients", r.ingredientController.CreateIngredient)
	app.Post("/recipes", r.recipeController.CreateRecipe)

}

func (r Router) Start(app *fiber.App) {
	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
