package e2etest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/denisyao1/welsh-academy-api/model"
	"github.com/denisyao1/welsh-academy-api/repository"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	t.Parallel()
	// create admin user if not exist
	userService.CreateIfNotExist(&model.User{Username: "admin", Password: "admin", IsAdmin: true})

	assert := assert.New(t)
	testCases := []struct {
		inputs      []byte
		statusCode  int
		description string
	}{
		{
			inputs:      []byte(`{"username":"admin", "password": "admin"}`),
			statusCode:  200,
			description: "login with correct credentials",
		},
		{
			inputs:      []byte(`{"username":"adm", "password": "admin"}`),
			statusCode:  401,
			description: "login with wrong user name",
		},
		{
			inputs:      []byte(`{"username":"ad", "password": "admin"}`),
			statusCode:  401,
			description: "login with user name length < 3",
		},
		{
			inputs:      []byte(`{"username":"admin", "password": "ad"}`),
			statusCode:  401,
			description: "login with wrong password",
		},
		{
			inputs:      []byte(`{"username":"admin", "password": "admin01"}`),
			statusCode:  401,
			description: "login with wrong password",
		},
	}
	url := "/api/v1/login"
	for _, tt := range testCases {
		req := httptest.NewRequest(PostMethod, url, bytes.NewBuffer(tt.inputs))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := App.Test(req, -1)
		assert.Equal(tt.statusCode, resp.StatusCode, tt.description)
		cookie := resp.Cookies()
		if resp.StatusCode == 200 {
			assert.Equal(1, len(cookie), tt.description)
		} else {
			assert.Equal(0, len(cookie), tt.description)
		}
	}

}
func TestLogout(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	// create admin user if not exist
	userService.CreateIfNotExist(&model.User{Username: "admin", Password: "admin", IsAdmin: true})

	// login first
	code, _ := login("admin", "admin")
	if code != 200 {
		t.Log("Auth failed")
		t.FailNow()
	}

	// Logout and try to access protecting ressources
	url := "/api/v1/logout"
	req := httptest.NewRequest(GetMethod, url, nil)
	resp, _ := App.Test(req, -1)
	assert.Equal(200, resp.StatusCode, "logout status code not OK")

	testCases := []struct {
		url         string
		statusCode  int
		description string
	}{
		{
			url:         "/api/v1/ingredients",
			statusCode:  401,
			description: "List ingredients unauthorized",
		},
		{
			url:         "/api/v1/users/my-infos",
			statusCode:  401,
			description: "Access my infos unauthorized",
		},
		{
			url:         "/api/v1/health",
			statusCode:  200,
			description: "Access health check authorized",
		},
	}

	for _, tt := range testCases {
		req := httptest.NewRequest(GetMethod, tt.url, nil)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := App.Test(req, -1)
		assert.Equal(tt.statusCode, resp.StatusCode, tt.description)
	}
}

func TestCreateUser(t *testing.T) {
	// t.Parallel()
	assert := assert.New(t)

	// create admin user if not exist
	userService.CreateIfNotExist(&model.User{Username: "admin", Password: "admin", IsAdmin: true})

	// login admin user first
	code, authCookie := login("admin", "admin")
	if code != 200 {
		t.Log("Auth failed")
		t.FailNow()
	}

	testCases := []struct {
		username    string
		password    string
		isAdmin     bool
		statusCode  int
		description string
	}{
		{
			username:    "te",
			password:    "test",
			statusCode:  400,
			description: "short username, status should be bad request",
		},
		{
			username:    "test",
			password:    "tes",
			statusCode:  400,
			description: "short paswword, status should be bad request",
		},
		{
			username:    "admin",
			password:    "test",
			statusCode:  409,
			description: "existing username, status should be conflict",
		},
		{
			username:    "testUser",
			password:    "testUser",
			statusCode:  201,
			description: "correct inputs, status shoud be Created",
		},
		{
			username:    "testUser",
			password:    "testUser",
			statusCode:  409,
			description: "existing username, status should be conflict",
		},
		{
			username:    "testAdmin",
			password:    "testAdmin",
			statusCode:  201,
			description: "good username and password, status should be Created",
			isAdmin:     true,
		},
	}

	url := "/api/v1/users"
	for _, tt := range testCases {
		json_inputs := fmt.Sprintf(`{"username":"%s", "password":"%s","isAdmin":"%v"}`,
			tt.username, tt.password, tt.isAdmin)
		inputs := []byte(json_inputs)
		req := httptest.NewRequest(PostMethod, url, bytes.NewBuffer(inputs))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(authCookie)
		resp, _ := App.Test(req, -1)
		assert.Equal(tt.statusCode, resp.StatusCode, tt.description)
		if resp.StatusCode == 200 {
			userRepo := repository.NewUserRepository(InMemoryDB.GetDB())
			user := model.User{Username: "test"}
			userRepo.GetByUsername(&user)
			assert.NotEqual(0, user.ID, "test user should be in the DB")
			if tt.isAdmin {
				assert.True(user.IsAdmin, "Admin user should have isAdmin set to true in DB")
			} else {
				assert.False(user.IsAdmin, "Non admin user should have isAdmin set to false in DB")
			}
		}
	}

}

func TestUpdatePassword(t *testing.T) {
	// t.Parallel()
	assert := assert.New(t)

	// if not exist create user test
	userService.CreateIfNotExist(&model.User{Username: "testPass", Password: "test", IsAdmin: false})

	// login
	code, authCookie := login("testPass", "test")
	if code != 200 {
		t.Log("Auth failed")
		t.FailNow()
	}

	// run test case using Auth cookie
	testCases := []struct {
		password    string
		statusCode  int
		description string
	}{
		{
			password:    "te",
			statusCode:  400,
			description: "short password, should return Bad Request",
		},
		{
			password:    "password",
			statusCode:  200,
			description: "good, should return OK",
		},
		{
			password:    "password",
			statusCode:  400,
			description: "same password has precedent, should return Bad Request",
		},
		{
			password:    "test",
			statusCode:  200,
			description: "new password, should return OK",
		},
	}
	url := "/api/v1/users/password-change"

	for _, tt := range testCases {
		json := fmt.Sprintf(`{"password":"%s"}`, tt.password)
		inputs := []byte(json)
		req := httptest.NewRequest(PatchMethod, url, bytes.NewBuffer(inputs))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(authCookie)
		resp, _ := App.Test(req, -1)
		assert.Equal(tt.statusCode, resp.StatusCode, tt.description)
		if resp.StatusCode == 200 {
			code, _ := login("testPass", tt.password)
			assert.Equal(200, code, "Login  should not failed after password update")
		}

	}

}

func TestGetUserInfos(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	// if not exist create user test
	userService.CreateIfNotExist(&model.User{Username: "test", Password: "test", IsAdmin: false})

	testCases := []struct {
		username    string
		password    string
		admin       bool
		statusCode  int
		description string
	}{
		{
			username:    "testify",
			password:    "test",
			statusCode:  401,
			description: "wrong credentials, should return unauthorized",
		},
		{
			username:    "test",
			password:    "test",
			admin:       false,
			statusCode:  200,
			description: "good credential, should be OK",
		},
		{
			username:    "admin",
			password:    "admin",
			admin:       true,
			statusCode:  200,
			description: "good credential, should be OK",
		},
	}

	url := BaseUrl + "/users/my-infos"

	for _, tt := range testCases {
		code, authCookie := login(tt.username, tt.password)
		req := httptest.NewRequest(GetMethod, url, nil)
		if code == 200 {
			req.AddCookie(authCookie)
		}
		resp, _ := App.Test(req, -1)
		msg, _ := io.ReadAll(resp.Body)
		assert.Equal(tt.statusCode, resp.StatusCode, tt.description)
		if resp.StatusCode == 200 {
			u := model.User{}
			json.Unmarshal(msg, &u)
			v := u.IsAdmin
			assert.Equal(tt.admin, v, "expected admin=%v but got value=%v", tt.admin, v)
		}

	}

}
