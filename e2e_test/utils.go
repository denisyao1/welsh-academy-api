package e2etest

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/denisyao1/welsh-academy-api/common"
	"github.com/denisyao1/welsh-academy-api/controller"
	"github.com/denisyao1/welsh-academy-api/database"
	"github.com/denisyao1/welsh-academy-api/repository"
	"github.com/denisyao1/welsh-academy-api/router"
	"github.com/denisyao1/welsh-academy-api/service"
	"github.com/gofiber/fiber/v2"
)

var (
	GetMethod   = "GET"
	PostMethod  = "POST"
	PatchMethod = "PATCH"
	BaseUrl     = "/api/v1"
)

var (
	BadRequest   = 400
	OK           = 200
	Created      = 201
	Unauthorized = 401
	Conflict     = 409
)

var (
	InMemoryDB     database.GormDB
	Config         = common.Configuration{JWT_SECRET: "test"}
	userRepo       repository.UserRepository
	userService    service.UserService
	ingredientRepo repository.IngredientRepository
	recipeRepo     repository.RecipeRepository
	App            = CreateTestApp()
)

func CreateTestApp() *fiber.App {
	var err error
	InMemoryDB, err = database.NewInMemoryDB(true)
	if err != nil {
		log.Fatalln("Unable te create In Memory Database")
	}

	// migrate database
	InMemoryDB.MigrateAll()

	ingredientRepo = repository.NewGormIngredientRepository(InMemoryDB.GetDB())
	ingredientService := service.NewIngredientService(ingredientRepo)
	ingredienController := controller.NewIngredientController(ingredientService)

	recipeRepo = repository.NewGormRecipeRepository(InMemoryDB.GetDB())
	recipeService := service.NewRecipeService(recipeRepo, ingredientRepo)
	recipeController := controller.NewRecipeController(recipeService)

	userRepo = repository.NewUserRepository(InMemoryDB.GetDB())

	userService = service.NewUserService(userRepo, Config.JWT_SECRET)

	// create default admin user
	userService.CreateDefaultAdmin()

	userController := controller.NewUserController(userService)

	router := router.New(ingredienController, recipeController, userController, Config.JWT_SECRET)

	app := fiber.New()

	router.InitRoutes(app)

	return app
}

func login(username, password string) (int, *http.Cookie) {

	json := fmt.Sprintf(`{"username":"%s", "password": "%s"}`, username, password)
	inputs := []byte(json)
	url := "/api/v1/login"
	req := httptest.NewRequest(PostMethod, url, bytes.NewBuffer(inputs))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := App.Test(req, -1)
	cookies := resp.Cookies()

	// test status code and cookie
	if len(cookies) == 0 {
		return resp.StatusCode, nil
	}
	return resp.StatusCode, cookies[0]

}
