package main

import (
	"log"

	"github.com/denisyao1/welsh-academy-api/controller"
	"github.com/denisyao1/welsh-academy-api/database"
	"github.com/denisyao1/welsh-academy-api/initializer"
	"github.com/denisyao1/welsh-academy-api/repository"
	"github.com/denisyao1/welsh-academy-api/router"
	"github.com/denisyao1/welsh-academy-api/service"
	"github.com/gofiber/fiber/v2"
)

func main() {
	envErr := initializer.LoadEnvVariables()
	if envErr != nil {
		log.Fatal("Failed to load env variable")
	}

	gormDB := database.NewGormDB()

	ingredientRepo := repository.NewGormIngredientRepository(gormDB.GetDB())
	ingredientService := service.NewIngredientService(ingredientRepo)
	ingredienController := controller.NewIngredientController(ingredientService)

	recipeRepo := repository.NewGormRecipeRepository(gormDB.GetDB())
	recipeService := service.NewRecipeService(recipeRepo, ingredientRepo)
	recipeController := controller.NewRecipeController(recipeService)

	userRepo := repository.NewUserRepository(gormDB.GetDB())
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	router := router.New(ingredienController, recipeController, userController)

	app := fiber.New()

	router.InitRoutes(app)

	router.Start(app)

}
