package controller

import "github.com/gofiber/fiber/v2"

//	HealthCheck shows if Api is running.
//
// @Summary      Health check
// @Description  Check Api is running
// @Tags         Health
// @Produce      json
// @Success      200 {object} Message
// @Failure      500
// @Router       /health [get]
func HealthCheck(c *fiber.Ctx) error {
	return c.JSON(NewMessage("Welsh Academy Api is running."))
}
