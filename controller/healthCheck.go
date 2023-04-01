package controller

import "github.com/gofiber/fiber/v2"

func HealthCheck(c *fiber.Ctx) error {
	return c.JSON(Map{"Message": "Welsh Academy Api is running."})
}
