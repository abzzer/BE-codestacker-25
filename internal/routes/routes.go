package routes

import (
	"github.com/abzzer/BE-codestacker-25/internal/middleware"
	"github.com/abzzer/BE-codestacker-25/internal/repository"
	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes sets up all the API routes
func RegisterRoutes(app *fiber.App) {
	// Public Route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("We have a working API with databases and RBAC!!")
	})

	// Register a user
	app.Post("/register", func(c *fiber.Ctx) error {
		type Request struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Role     string `json:"role"`
		}

		var req Request
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
		}

		user, err := repository.CreateUser(req.Username, req.Password, req.Role)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Could not create user"})
		}

		return c.JSON(user)
	})

	// Protected Route (Admin Only)
	app.Get("/admin", middleware.AuthMiddleware("admin"), func(c *fiber.Ctx) error {
		return c.SendString("Welcome Admin!")
	})
}
