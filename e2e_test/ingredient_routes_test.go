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

func TestIngredientCreate(t *testing.T) {
	// t.Parallel()
	assert := assert.New(t)

	// create admin user if not exist
	userService.CreateIfNotExist(&model.User{Username: "admin", Password: "admin", IsAdmin: true})

	// login admin
	code, authCookie := login("admin", "admin")
	if code != 200 {
		t.Log("Auth failed")
		t.FailNow()
	}

	testCases := []struct {
		name        string
		statusCode  int
		numberInDB  int
		description string
	}{
		{
			statusCode:  400,
			description: "Empty name, should be bad request",
		},
		{
			name:        "ingredient1",
			statusCode:  201,
			numberInDB:  1,
			description: "valid name, should be created",
		},
		{
			name:        "ingredient1",
			statusCode:  409,
			numberInDB:  1,
			description: "alredy used name, should return conflict",
		},
		{
			name:        "ingredient2",
			statusCode:  201,
			numberInDB:  2,
			description: "valid name, should return created",
		},
	}

	url := BaseUrl + "/ingredients"

	for _, tt := range testCases {
		json_inputs := fmt.Sprintf(`{"name":"%s"}`, tt.name)
		inputs := []byte(json_inputs)
		req := httptest.NewRequest(PostMethod, url, bytes.NewBuffer(inputs))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(authCookie)
		resp, _ := App.Test(req, -1)
		assert.Equal(tt.statusCode, resp.StatusCode, tt.description)
		if resp.StatusCode == 201 {
			ingredients, _ := ingredientRepo.FindAll()
			v := len(ingredients)
			assert.Equal(tt.numberInDB, v, "expect %v ingredients in DB got %v", tt.numberInDB, v)
		}
	}

}

func TestIngredientsListing(t *testing.T) {
	// t.Parallel()
	assert := assert.New(t)

	// create two users if not exists
	userService.CreateIfNotExist(&model.User{Username: "admin", Password: "admin"})
	userService.CreateIfNotExist(&model.User{Username: "test", Password: "test"})

	// test cases
	testCases := []struct {
		username    string
		password    string
		statusCode  int
		description string
	}{
		{
			username:    "test2",
			password:    "test2",
			statusCode:  401,
			description: "invalid credential, should be unauthorized",
		},
		{
			username:    "test",
			password:    "test",
			statusCode:  200,
			description: "valid credential, should be Ok",
		},
		{
			username:    "admin",
			password:    "admin",
			statusCode:  200,
			description: "valid credential, should be Ok",
		},
	}

	url := BaseUrl + "/ingredients"
	for _, tt := range testCases {
		code, authCookie := login(tt.username, tt.password)
		assert.Equal(tt.statusCode, code, tt.description)
		if code != 200 {
			continue
		}
		req := httptest.NewRequest(GetMethod, url, nil)
		req.AddCookie(authCookie)
		resp, _ := App.Test(req, -1)
		results, _ := io.ReadAll(resp.Body)
		response := schema.IngredientsResponse{}
		json.Unmarshal(results, &response)
		v := len(response.Ingredients)
		assert.Equal(2, v, "request response should show 3 ingredients but got %v", v)
	}

}
