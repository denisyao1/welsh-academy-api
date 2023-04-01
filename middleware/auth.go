package middleware

import (
	"fmt"
	"time"

	"github.com/denisyao1/welsh-academy-api/model"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func JwtWare(siginKey string, role model.Role) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		// @TODO avoid non defined route
		tokenString := ctx.Cookies("jwt")
		if tokenString == "" {
			return ctx.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"message": "Missing or malformed JWT"})
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(siginKey), nil
		})

		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"message": "Missing or malformed JWT"})
		}

		// reference ctx json missing token response
		invalidTokenResponse := ctx.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"message": "invalid or expired token"})

		claims, ok := token.Claims.(jwt.MapClaims)
		if !(ok && token.Valid) {
			return invalidTokenResponse
		}

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return invalidTokenResponse
		}

		user_role := model.Role(claims["role"].(float64))
		if role == model.RoleAdmin && role != user_role {
			return invalidTokenResponse
		}
		ctx.Locals("userID", claims["ID"])
		return ctx.Next()
	}
}
