package controller

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

const (
	BadRequest   = fiber.StatusBadRequest
	Conflict     = fiber.StatusConflict
	Created      = fiber.StatusCreated
	OK           = fiber.StatusOK
	Unauthorized = fiber.StatusUnauthorized
)

type Map = fiber.Map

type BaseController struct{}

func (b BaseController) HandleUnExpetedError(err error, ctx *fiber.Ctx) error {
	log.Println("UnExpectedError: ", err.Error())
	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
}

func (b BaseController) GetConnectedUserID(ctx *fiber.Ctx) (int, error) {
	return strconv.Atoi(ctx.Locals("userID").(string))
}

func (b BaseController) ConvertParamToInt(paramName string, ctx *fiber.Ctx) (int, error) {
	return strconv.Atoi(ctx.Params(paramName))
}
