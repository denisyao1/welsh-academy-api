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
	service service.UserService
}

func NewUserController(service service.UserService) UserController {
	fmt.Println("service=", service)
	return UserController{service: service}
}

func (c UserController) CreateUser(ctx *fiber.Ctx) error {
	var userSchema schema.CreateUserSchema

	if err := ctx.BodyParser(&userSchema); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to read request body"})
	}
	validationErrs := c.service.ValidateUserCreation(userSchema)

	if validationErrs != nil {
		if len(validationErrs) == 1 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": validationErrs[0]})
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": validationErrs})
	}

	user, err := c.service.CreateUser(userSchema)

	if err != nil {
		if errors.Is(err, exception.ErrDuplicateKey) {
			message := fmt.Sprintf("username '%s' already exists.", user.Username)
			return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{"error": message})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})

	}

	return ctx.Status(fiber.StatusCreated).JSON(user)
}

func (c UserController) Login(ctx *fiber.Ctx) error {
	// get request body
	var loginSchema schema.LoginSchema
	if err := ctx.BodyParser(&loginSchema); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to read request body"})
	}

	// create token
	token, err := c.service.CreateAccessToken(loginSchema)
	if err != nil {
		if errors.Is(err, exception.ErrInvalidCredentials) {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
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

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "login successful"})
}

func (c UserController) Logout(ctx *fiber.Ctx) error {
	ctx.ClearCookie("jwt")
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "logout successful"})
}

// @TODO: I have to complete or remove this fonction
func (c UserController) UserInfos(ctx *fiber.Ctx) error {
	userID, err := strconv.Atoi(ctx.Locals("userID").(string))
	if err != nil {
		ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "invalid or expired token"})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"userID": userID})
}
