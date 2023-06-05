package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(ctx *fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Missing authorization header",
		})
	}

	tokenString := ""
	if len(authHeader) > 7 {
		tokenString = authHeader[7:]
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.ErrUnauthorized
		}
		return []byte("secret-key"), nil
	})
	if err != nil || !token.Valid {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid or expired token",
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid token claims",
		})
	}

	role := claims["role"].(string)

	ctx.Locals("role", role)
	return ctx.Next()
}
