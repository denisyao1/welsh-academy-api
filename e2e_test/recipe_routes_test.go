package e2etest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/denisyao1/welsh-academy-api/model"
	"github.com/denisyao1/welsh-academy-api/schema"
	"github.com/stretchr/testify/assert"
)

func TestCreateRecipe(t *testing.T) {
	// t.Parallel()

	assert := assert.New(t)

	// create an admin user if not exists
	userService.CreateIfNotExist(&model.User{Username: "admin", Password: "admin", IsAdmin: true})

	// create some ingredients if not exist
	ingredientRepo.GetOrCreate("ingredient01")
	ingredientRepo.GetOrCreate("ingredient02")

	// login the admin user
	code, authCookie := login("admin", "admin")
	if code != 200 {
		t.Log("Auth faild")
		t.FailNow()
	}

	// testCases
	testCases := []struct {
		name        string
		making      string
		ingredients string
		statusCode  int
		description string
	}{

		{
			name:        "recipe0",
			ingredients: `[{"name":"ingredient01"}]`,
			statusCode:  BadRequest,
			description: "making is empty should return BadRequest",
		},
		{
			name:        "recipe0",
			making:      "Mix all",
			ingredients: `[]`,
			statusCode:  BadRequest,
			description: "ingredients is empty should return BadRequest",
		},
		{
			name:        "recipe0",
			making:      "Mix all",
			ingredients: `[{"name":"Onion"},{"name":"Tomato"}]`,
			statusCode:  BadRequest,
			description: "Contains undefined ingredients should return BadRequest",
		},
		{
			name:        "recipe0",
			making:      "Mix all",
			ingredients: `[{"name":"ingredient01"},{"name":"Tomato"}]`,
			statusCode:  BadRequest,
			description: "Contains one undefined ingredient should return BadRequest",
		},
		{
			name:        "recipe0",
			making:      "Mix all",
			ingredients: `[{"name":"ingredient01"},{"name":"ingredient02"}]`,
			statusCode:  Created,
			description: "Valid input, recipe should be created",
		},
	}

	url := BaseUrl + "/recipes"

	for _, tt := range testCases {
		json_inputs := fmt.Sprintf(`{"name":"%s", "making":"%s", "ingredients":%s}`, tt.name, tt.making, tt.ingredients)
		inputs := []byte(json_inputs)
		req := httptest.NewRequest(PostMethod, url, bytes.NewBuffer(inputs))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(authCookie)
		resp, _ := App.Test(req, -1)
		assert.Equal(tt.statusCode, resp.StatusCode, tt.description)
		if resp.StatusCode != Created {
			continue
		}
		result, _ := io.ReadAll(resp.Body)
		var recipe model.Recipe
		json.Unmarshal(result, &recipe)
		name := recipe.Name
		making := recipe.Making
		l := len(recipe.Ingredients)
		assert.Equal(tt.name, name, "recipe name should be %s but got %s", tt.name, name)
		assert.Equal(tt.making, making, "recipe name should be %s but got %s", tt.making, making)
		assert.Equal(2, l, "recipe should have %v ingredients but got %v", 2, l)
	}
}

func TestListAllPossibleRecipes(t *testing.T) {
	// t.Parallel()

	assert := assert.New(t)

	//create some ingredients
	ingredient1, _ := ingredientRepo.GetOrCreate("ingredient1")
	ingredient2, _ := ingredientRepo.GetOrCreate("ingredient2")
	ingredient3, _ := ingredientRepo.GetOrCreate("ingredient3")
	ingredientA, _ := ingredientRepo.GetOrCreate("ingredientA")
	// create some recipes
	recipe1 := model.Recipe{
		Name:        "recipe1",
		Making:      "making recip2",
		Ingredients: []model.Ingredient{ingredient1, ingredient2}}
	recipe2 := model.Recipe{
		Name:        "recipe2",
		Making:      "dummy",
		Ingredients: []model.Ingredient{ingredient2, ingredientA}}

	recipe3 := model.Recipe{
		Name:        "recipe3",
		Making:      "dummy",
		Ingredients: []model.Ingredient{ingredient3}}

	recipeRepo.GetOrCreate(&recipe1)
	recipeRepo.GetOrCreate(&recipe2)
	recipeRepo.GetOrCreate(&recipe3)

	recipes, _ := recipeRepo.FindAll()
	if len(recipes) < 3 {
		t.Log("it should be 3 or more recipes in the DB")
		t.FailNow()
	}

	// create admin user if not exist
	user := model.User{Username: "test", Password: "test", IsAdmin: false}
	userService.CreateIfNotExist(&user)

	// login
	code, authCookie := login("test", "test")
	if code != 200 {
		t.Log("Auth failed")
		t.FailNow()
	}

	// test cases
	testCases := []struct {
		ingredients string
		number      int
		description string
	}{
		{
			number:      3,
			description: "no filter on ingredients, should return at least 3 recipes but got %v",
		},
		{
			ingredients: "ingredient1",
			number:      1,
			description: "filter on ingredient1 should return 2 recipes but got %v",
		},
		{
			ingredients: "ingredient2",
			number:      2,
			description: "filter on ingredient2 should return 2 recipes but got %v",
		},
		{
			ingredients: "ingredient1,ingredient2",
			number:      2,
			description: "filter on ingredient 1 & 2 should return 2 recipes but got %v",
		},
		{
			ingredients: "ingredient1,ingredient2,ingredient3",
			number:      3,
			description: "filter on ingredient 1&2&3 should return 3 recipes but got %v",
		},
	}

	url := BaseUrl + "/recipes?ingredients="

	for _, tt := range testCases {
		req := httptest.NewRequest(GetMethod, url+tt.ingredients, nil)
		req.AddCookie(authCookie)
		resp, _ := App.Test(req, -1)
		results, _ := io.ReadAll(resp.Body)
		response := schema.RecipesResponse{}
		json.Unmarshal(results, &response)
		v := response.Count

		if tt.ingredients == "" {
			assert.True(v >= tt.number, tt.description, v)
		} else {
			assert.Equal(tt.number, v, tt.description, v)
		}

	}

}

func TestFlagUnflagRecipes(t *testing.T) {
	assert := assert.New(t)

	// create ingredients and recipes
	ingredient1, _ := ingredientRepo.GetOrCreate("ingredient1")
	ingredient2, _ := ingredientRepo.GetOrCreate("ingredient2")
	ingredientA, _ := ingredientRepo.GetOrCreate("ingredientA")

	recipe1 := model.Recipe{
		Name:        "recipe1",
		Making:      "making recip2",
		Ingredients: []model.Ingredient{ingredient1, ingredient2}}
	recipe2 := model.Recipe{
		Name:        "recipe2",
		Making:      "dummy",
		Ingredients: []model.Ingredient{ingredient2, ingredientA}}

	recipeRepo.GetOrCreate(&recipe1)
	recipeRepo.GetOrCreate(&recipe2)

	// get or create user and login
	user := model.User{Username: "test", Password: "test", IsAdmin: false}
	userService.CreateIfNotExist(&user)

	code, authCookie := login("test", "test")
	if code != 200 {
		t.Log("Auth failed")
		t.FailNow()
	}

	// test cases
	testCases := []struct {
		recipeID    int
		isPresent   bool
		number      int // number of recipe in user favorites
		description string
	}{
		{
			recipeID:    recipe1.ID,
			isPresent:   true,
			number:      1,
			description: "recipe 1 should be in user favorites",
		},
		{
			recipeID:    recipe2.ID,
			isPresent:   true,
			number:      2,
			description: "recipe 2 should be in user favorites",
		},
		{
			recipeID:    recipe1.ID,
			number:      1,
			description: "recipe 1 shouldn't be in user favorites",
		},
		{
			recipeID:    recipe2.ID,
			number:      0,
			description: "recipe 2 shouldn't be in user favorites",
		},
	}

	for _, tt := range testCases {
		url := fmt.Sprintf("%s/recipes/%v/flag-unflag", BaseUrl, tt.recipeID)
		req := httptest.NewRequest(PostMethod, url, nil)
		req.AddCookie(authCookie)
		resp, _ := App.Test(req, -1)
		assert.Equal(200, resp.StatusCode, "request should be OK")
		ok, _ := recipeRepo.IsInUserFavorites(user.ID, tt.recipeID)
		assert.True(ok == tt.isPresent, tt.description)
		recipes, _ := recipeRepo.FindFavorites(user.ID)
		l := len(recipes)
		assert.Equal(tt.number, l, "user should have %v recipes in its favorites but got %v", tt.number, l)
	}

}
