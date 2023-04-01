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

type UserController struct {
	BaseController
	service service.UserService
}

func NewUserController(service service.UserService) UserController {
	return UserController{service: service}
}

func (c UserController) Create(ctx *fiber.Ctx) error {
	var userSchema schema.CreateUserSchema

	if err := ctx.BodyParser(&userSchema); err != nil {
		return ctx.Status(BadRequest).JSON(Map{"error": "Failed to read request body"})
	}
	validationErrs := c.service.ValidateUserCreation(userSchema)

	if validationErrs != nil {
		if len(validationErrs) == 1 {
			return ctx.Status(BadRequest).JSON(Map{"error": validationErrs[0]})
		}
		return ctx.Status(BadRequest).JSON(Map{"errors": validationErrs})
	}

	user, err := c.service.CreateUser(userSchema)

	if err != nil {
		if errors.Is(err, exception.ErrDuplicateKey) {
			message := fmt.Sprintf("username '%s' already exists.", user.Username)
			return ctx.Status(Conflict).JSON(Map{"error": message})
		}
		return c.HandleUnExpetedError(err, ctx)
	}

	return ctx.Status(fiber.StatusCreated).JSON(user)
}

func (c UserController) Login(ctx *fiber.Ctx) error {
	// get request body
	var loginSchema schema.LoginSchema
	if err := ctx.BodyParser(&loginSchema); err != nil {
		return ctx.Status(BadRequest).JSON(Map{"error": "Failed to read request body"})
	}

	// create token
	token, err := c.service.CreateAccessToken(loginSchema)
	if err != nil {
		return c.errorHandlerHelper(err, ctx)
	}

	// create cookie
	jwt_cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		SameSite: "lax",
	}

	// set cookie
	ctx.Cookie(&jwt_cookie)

	return ctx.Status(OK).JSON(Map{"message": "login successful"})
}

func (c UserController) Logout(ctx *fiber.Ctx) error {
	ctx.ClearCookie("jwt")
	return ctx.Status(OK).JSON(Map{"message": "logout successful"})
}

func (c UserController) GetInfos(ctx *fiber.Ctx) error {
	userID, err := strconv.Atoi(ctx.Locals("userID").(string))
	if err != nil {
		ctx.Status(Unauthorized).JSON(Map{"message": "Missing or malformed JWT"})
	}

	user, err := c.service.GetInfos(userID)
	if err != nil {
		c.errorHandlerHelper(err, ctx)
	}

	return ctx.Status(fiber.StatusOK).JSON(user)
}

func (c UserController) UpdatePassword(ctx *fiber.Ctx) error {
	//get user id
	userID, err := c.GetConnectedUserID(ctx)
	if err != nil {
		return ctx.Status(Unauthorized).
			JSON(Map{"message": "Missing or malformed JWT"})
	}

	var pwdSchema schema.PasswordSchema

	// get password schema
	if err = ctx.BodyParser(&pwdSchema); err != nil {
		return ctx.Status(BadRequest).JSON(Map{"error": "can't read request body"})
	}

	//update password
	if err = c.service.UpdatePaswword(userID, pwdSchema); err != nil {
		return c.errorHandlerHelper(err, ctx)
	}

	return ctx.Status(OK).JSON(Map{"message": "password update successful"})
}

func (c UserController) errorHandlerHelper(err error, ctx *fiber.Ctx) error {
	if errors.Is(err, exception.ErrInvalidPassword) || errors.Is(err, exception.ErrPasswordSame) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if errors.Is(err, exception.ErrRecordNotFound) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user " + err.Error()})
	}

	if errors.Is(err, exception.ErrInvalidCredentials) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}

	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
}
