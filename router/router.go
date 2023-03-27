package router

import (
	"log"
	"os"

	"github.com/denisyao1/welsh-academy-api/controllers"
	"github.com/gofiber/fiber/v2"
)

type Router struct {
	ingredientController controllers.IngredientController
}

func New(ingredientController controllers.IngredientController) *Router {
	return &Router{ingredientController: ingredientController}
}

func (r Router) InitRoutes(app *fiber.App) {
	app.Get("/health", controllers.HealthCheck)
	app.Post("/ingredients", r.ingredientController.CreateIngredient)
	app.Get("/ingredients", r.ingredientController.ListAllIngredients)

}

func (r Router) Start(app *fiber.App) {
	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
