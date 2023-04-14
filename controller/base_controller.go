package controller

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// BaseController contains common method for all controllers.
type BaseController struct{}

// HandleUnExpetedError handles errors the api didn't except.
func (b BaseController) HandleUnExpetedError(err error, ctx *fiber.Ctx) error {
	log.Println("UnExpectedError: ", err.Error())
	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
}

// GetConnectedUserID returns the connected user id.
func (b BaseController) GetConnectedUserID(ctx *fiber.Ctx) (int, error) {
	return strconv.Atoi(ctx.Locals("userID").(string))
}

// ConvertParamToInt convert path or query param to int.
func (b BaseController) ConvertParamToInt(paramName string, ctx *fiber.Ctx) (int, error) {
	return strconv.Atoi(ctx.Params(paramName))
}
