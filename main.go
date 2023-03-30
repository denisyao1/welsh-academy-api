package main

import (
	"log"

	"github.com/denisyao1/welsh-academy-api/controllers"
	"github.com/denisyao1/welsh-academy-api/database"
	"github.com/denisyao1/welsh-academy-api/initializers"
	"github.com/denisyao1/welsh-academy-api/repositories"
	"github.com/denisyao1/welsh-academy-api/router"
	"github.com/denisyao1/welsh-academy-api/services"
	"github.com/gofiber/fiber/v2"
)

func main() {
	envErr := initializers.LoadEnvVariables()
	if envErr != nil {
		log.Fatal("Failed to load env variable")
	}

	gormDB := database.NewGormDB()

	ingredientRepo := repositories.NewGormIngredientRepository(gormDB.GetDB())
	ingredientService := services.NewIngredientService(ingredientRepo)
	ingredienController := controllers.NewIngredientController(ingredientService)

	recipeRepo := repositories.NewGormRecipeRepository(gormDB.GetDB())
	recipeService := services.NewRecipeService(recipeRepo, ingredientRepo)
	recipeController := controllers.NewRecipeController(recipeService)

	userRepo := repositories.NewUserRepository(gormDB.GetDB())
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	router := router.New(ingredienController, recipeController, userController)

	app := fiber.New()

	router.InitRoutes(app)

	router.Start(app)

}
