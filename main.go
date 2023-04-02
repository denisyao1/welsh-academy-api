package main

import (
	"log"

	"github.com/denisyao1/welsh-academy-api/common"
	"github.com/denisyao1/welsh-academy-api/controller"
	"github.com/denisyao1/welsh-academy-api/database"
	"github.com/denisyao1/welsh-academy-api/repository"
	"github.com/denisyao1/welsh-academy-api/router"
	"github.com/denisyao1/welsh-academy-api/schema"
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
	userController := controller.NewUserController(userService)

	// create an default admin user
	var user schema.CreateUserSchema
	user.Username = "admin"
	user.Password = "admin"
	user.IsAdmin = true
	_, err := userService.CreateUser(user)
	if err != nil {
		log.Fatalln("Failed to create default admin user")
	}
	log.Println(" Default admin user created")

	router := router.New(ingredienController, recipeController, userController, config.JWT_SECRET)

	app := fiber.New()

	app.Use(logger.New())

	router.InitRoutes(app)

	router.Start(app)

}
