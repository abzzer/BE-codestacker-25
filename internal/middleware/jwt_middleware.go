package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// Accepts a list of allowed roles
func JWTMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing Authorization header"})
		}

		parts := strings.Split(authHeader, "Bearer ")
		if len(parts) != 2 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid Authorization format"})
		}

		tokenStr := parts[1]
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "supersecretkey"
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
		}

		role, ok := claims["role"].(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid role claim"})
		}

		// Check if the role is in allowedRoles
		allowed := false
		for _, r := range allowedRoles {
			if role == r {
				allowed = true
				break
			}
		}

		if !allowed {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
		}

		c.Locals("user_id", claims["user_id"])
		c.Locals("role", claims["role"])

		return c.Next()
	}
}
