package controller

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/denisyao1/welsh-academy-api/exception"
	"github.com/denisyao1/welsh-academy-api/schema"
	"github.com/denisyao1/welsh-academy-api/service"
	"github.com/gofiber/fiber/v2"
)

// User controller contains informations and methods to route user related requests.
type UserController struct {
	BaseController
	service service.UserService
}

// NewUserController creates new user controller
func NewUserController(service service.UserService) UserController {
	return UserController{service: service}
}

//	Create creates new user
//
// @Summary      Create user
// @Description  Create user.
// @Description
// @Description  Require Admin Role.
// @Param request body schema.User true "User object"
// @Tags         User Management
// @Accept       json
// @Produce      json
// @Success      201 {object} model.User
// @Failure      400 {object} ErrMessage
// @Failure      401 {object} ErrMessage
// @Failure      500
// @Security JWT
// @Router       /users [post]
func (c UserController) Create(ctx *fiber.Ctx) error {
	var userSchema schema.User

	if err := ctx.BodyParser(&userSchema); err != nil {
		return ctx.Status(BadRequest).JSON(NewErrMessage("Failed to read request body"))
	}

	validationErrs := c.service.ValidateUserCreation(userSchema)
	if validationErrs != nil {
		if len(validationErrs) == 1 {
			return ctx.Status(BadRequest).JSON(Map{"error": validationErrs[0]})
		}
		return ctx.Status(BadRequest).JSON(Map{"errors": validationErrs})
	}

	user, err := c.service.Create(userSchema)

	if err != nil {
		if errors.Is(err, exception.ErrDuplicateKey) {
			message := fmt.Sprintf("username '%s' already exists.", user.Username)
			return ctx.Status(Conflict).JSON(Map{"error": message})
		}
		return c.HandleUnExpetedError(err, ctx)
	}

	return ctx.Status(fiber.StatusCreated).JSON(user)
}

//	Login creates new cookie access token
//
// @Summary      Login
// @Description  Get new cookie access token
// @Param request body schema.Login true "Credentials"
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200 {object} Message
// @Failure      401 {object} ErrMessage
// @Failure      500
// @Router       /login [post]
func (c UserController) Login(ctx *fiber.Ctx) error {
	// get request body
	var loginSchema schema.Login
	if err := ctx.BodyParser(&loginSchema); err != nil {
		return ctx.Status(BadRequest).JSON(NewErrMessage("Failed to read request body."))
	}

	// create token
	token, err := c.service.CreateAccessToken(loginSchema)
	if err != nil {
		//if err is ErrInvalidCredentials
		if errors.Is(err, exception.ErrInvalidCredentials) {
			return ctx.Status(Unauthorized).JSON(NewErrMessage("invalid credentials"))
		}
		return c.HandleUnExpetedError(err, ctx)
	}

	// create cookie
	jwt_cookie := fiber.Cookie{
		Name:     "Auth",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		SameSite: "lax",
	}

	// set cookie
	ctx.Cookie(&jwt_cookie)
	return ctx.Status(OK).JSON(Map{"message": "login successful"})
}

//	Logout Logs out the connected user
//
// @Summary      Logout
// @Description  Logout
// @Tags         Auth
// @Produce      json
// @Success      200 {object} Message
// @Router       /logout [get]
func (c UserController) Logout(ctx *fiber.Ctx) error {
	// create an expired cookie
	cookie := fiber.Cookie{
		Name:     "Auth",
		Value:    "deleted",
		Expires:  time.Now().Add(-5 * time.Second),
		HTTPOnly: true,
		SameSite: "lax",
	}

	ctx.Cookie(&cookie)

	return ctx.Status(OK).JSON(Map{"message": "logout successful"})
}

//	GetInfos returns connected user informations
//
// @Summary      My infos
// @Description  Show connected user informations.
// @Tags         User Profile
// @Accept       json
// @Produce      json
// @Success      200 {object} model.User
// @Failure      400 {object} ErrMessage
// @Failure      401 {object} ErrMessage
// @Failure      404 {object} ErrMessage
// @Failure      409 {object} ErrMessage
// @Failure      500
// @Security JWT
// @Router       /users/my-infos [get]
func (c UserController) GetInfos(ctx *fiber.Ctx) error {
	userID, err := strconv.Atoi(ctx.Locals("userID").(string))
	if err != nil {
		ctx.Status(Unauthorized).JSON(Map{"message": "Missing or malformed token"})
	}
	user, err := c.service.GetInfos(userID)
	if err != nil {
		if errors.Is(err, exception.ErrRecordNotFound) {
			return ctx.Status(NotFound).JSON(Map{"error": "user " + err.Error()})
		}
		c.HandleUnExpetedError(err, ctx)
	}

	return ctx.Status(fiber.StatusOK).JSON(user)
}

//	UpdatePassword updates connected user's password
//
// @Summary      Update password
// @Description  Update connected user's password
// @Param request body schema.Password true "Password"
// @Tags         User Profile
// @Accept       json
// @Produce      json
// @Success      200 {object} Message
// @Failure      400 {object} ErrMessage
// @Failure      401 {object} ErrMessage
// @Failure      404 {object} ErrMessage
// @Failure      500
// @Security JWT
// @Router       /users/password-change [patch]
func (c UserController) UpdatePassword(ctx *fiber.Ctx) error {
	//get user id
	userID, err := c.GetConnectedUserID(ctx)
	if err != nil {
		return ctx.Status(Unauthorized).
			JSON(NewErrMessage("Missing or malformed token"))
	}

	// get password input
	var pwdSchema schema.Password
	if err = ctx.BodyParser(&pwdSchema); err != nil {
		return ctx.Status(BadRequest).JSON(NewErrMessage("can't read request body"))
	}

	//update password
	if err = c.service.UpdatePaswword(userID, pwdSchema); err != nil {
		if errors.Is(err, exception.ErrInvalidPassword) ||
			errors.Is(err, exception.ErrPasswordSame) {
			return ctx.Status(BadRequest).JSON(NewErrMessage(err.Error()))
		}
		if errors.Is(err, exception.ErrRecordNotFound) {
			msg := "user " + err.Error()
			return ctx.Status(NotFound).JSON(NewErrMessage(msg))
		}
		return c.HandleUnExpetedError(err, ctx)
	}

	return ctx.Status(OK).JSON(NewMessage("password update successful"))
}
