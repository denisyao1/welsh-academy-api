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
	key := r.SigningKey
	user := model.RoleUser
	admin := model.RoleAdmin
	jware := middleware.JwtWare

	api := app.Group("/api/v1")

	// routes thant required no auth
	api.Get("/health", controller.HealthCheck)
	api.Post("/login", r.userController.Login)
	api.Get("/logout", r.userController.Logout)

	// required user auth routes
	api.Get("/ingredients", jware(key, user), r.ingredientController.ListIngredients)
	api.Get("/recipes", jware(key, user), r.recipeController.ListRecipes)
	api.Post("/recipes/:id/flag-unflag", jware(key, user), r.recipeController.FlagOrUnflag)
	api.Get("/recipes/favorites", jware(key, user), r.recipeController.ListUserFavorites)
	api.Get("/users/my-infos", jware(key, user), r.userController.GetInfos)
	api.Patch("/users/password-change", jware(key, user), r.userController.UpdatePassword)

	// required admin auth routes
	api.Post("/users", jware(key, admin), r.userController.Create)
	api.Post("/ingredients", jware(key, admin), r.ingredientController.CreateIngredient)
	api.Post("/recipes", jware(key, admin), r.recipeController.CreateRecipe)

}

func (r Router) Start(app *fiber.App) {
	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
