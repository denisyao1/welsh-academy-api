package main

import (
	"github.com/denisyao1/welsh-academy-api/common"
	"github.com/denisyao1/welsh-academy-api/controller"
	"github.com/denisyao1/welsh-academy-api/database"
	_ "github.com/denisyao1/welsh-academy-api/docs"
	"github.com/denisyao1/welsh-academy-api/repository"
	"github.com/denisyao1/welsh-academy-api/router"
	"github.com/denisyao1/welsh-academy-api/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
)

// @title Welsh Academy API
// @version 1.0
// @description Welsh Academy API

// @contact.name Denis YAO
// @contact.email denisyao@outlook.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /api/v1
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

	userService := service.NewUserService(userRepo, config.JWT_SECRET)

	// create default admin user
	userService.CreateDefaultAdmin()

	userController := controller.NewUserController(userService)

	router := router.New(ingredienController, recipeController, userController, config.JWT_SECRET)

	app := fiber.New()

	app.Use(logger.New())

	// swagger route
	app.Get("/docs/*", swagger.HandlerDefault)

	router.InitRoutes(app)

	router.Start(app)

}
