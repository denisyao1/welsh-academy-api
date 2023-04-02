package main

import (
	"github.com/denisyao1/welsh-academy-api/common"
	"github.com/denisyao1/welsh-academy-api/controller"
	"github.com/denisyao1/welsh-academy-api/database"
	"github.com/denisyao1/welsh-academy-api/repository"
	"github.com/denisyao1/welsh-academy-api/router"
	"github.com/denisyao1/welsh-academy-api/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	config := common.LoadConfig()

	gormDB := database.NewGormDB(config)

	// migrate database
	gormDB.Migrate()

	ingredientRepo := repository.NewGormIngredientRepository(gormDB.GetDB())
	ingredientService := service.NewIngredientService(ingredientRepo)
	ingredienController := controller.NewIngredientController(ingredientService)

	recipeRepo := repository.NewGormRecipeRepository(gormDB.GetDB())
	recipeService := service.NewRecipeService(recipeRepo, ingredientRepo)
	recipeController := controller.NewRecipeController(recipeService)

	userRepo := repository.NewUserRepository(gormDB.GetDB())

	userService := service.NewUserService(userRepo)

	// create an default admin user
	userService.CreateDefaultAdmin()

	userController := controller.NewUserController(userService)

	router := router.New(ingredienController, recipeController, userController, config.JWT_SECRET)

	app := fiber.New()

	app.Use(logger.New())

	router.InitRoutes(app)

	router.Start(app)

}
