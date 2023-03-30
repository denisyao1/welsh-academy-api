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
	app.Post("/ingredients", r.ingredientController.CreateIngredient)
	app.Get("/ingredients", r.ingredientController.ListAllIngredients)
	app.Post("/recipes", r.recipeController.CreateRecipe)
	app.Get("/recipes", r.recipeController.ListRecipes)

	app.Post("/users", r.userController.CreateUser)
	app.Post("/login", r.userController.Login)
	app.Get("/logout", r.userController.Logout)

	app.Use(middleware.JwtWare(r.SigningKey, model.RoleAdmin))
	app.Get("/infos", r.userController.UserInfos)

}

func (r Router) Start(app *fiber.App) {
	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
