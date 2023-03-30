package router

import (
	"log"
	"os"

	"github.com/denisyao1/welsh-academy-api/controllers"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

type Router struct {
	ingredientController controllers.IngredientController
	recipeController     controllers.RecipeController
	userController       controllers.UserController
}

func New(
	ingredientController controllers.IngredientController,
	recipeController controllers.RecipeController,
	userController controllers.UserController,
) *Router {
	return &Router{
		ingredientController: ingredientController,
		recipeController:     recipeController,
		userController:       userController,
	}
}

func (r Router) InitRoutes(app *fiber.App) {
	app.Get("/health", controllers.HealthCheck)
	app.Post("/ingredients", r.ingredientController.CreateIngredient)
	app.Get("/ingredients", r.ingredientController.ListAllIngredients)
	app.Post("/recipes", r.recipeController.CreateRecipe)
	app.Get("/recipes", r.recipeController.ListRecipes)

	app.Post("/users", r.userController.CreateUser)
	app.Post("/login", r.userController.Login)
	app.Get("/logout", r.userController.Logout)

	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}))

	// Protected routes
	app.Get("/infos", r.userController.UserInfos)
}

func (r Router) Start(app *fiber.App) {
	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
