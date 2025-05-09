package handlers

import (
	"github.com/abzzer/BE-codestacker-25/internal/auth"
	"github.com/abzzer/BE-codestacker-25/internal/models"
	"github.com/abzzer/BE-codestacker-25/internal/repository"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	passwordHash, role, clearance, err := repository.GetPasswordRoleClearanceByUserID(req.UserID)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "the repo not work",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user ID or password",
		})
	}

	token, err := auth.GenerateJWT(req.UserID, role, clearance)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not generate token",
		})
	}

	return c.JSON(fiber.Map{
		"token":  token,
		"userID": req.UserID,
		"role":   role,
	})
}

func LogoutHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Successfully logged out. Please discard your token on the client side.",
	})
}
