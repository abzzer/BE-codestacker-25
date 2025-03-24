package handlers

import (
	"github.com/abzzer/BE-codestacker-25/internal/models"
	"github.com/abzzer/BE-codestacker-25/internal/repository"
	"github.com/gofiber/fiber/v2"
)

func CreateUserHandler(c *fiber.Ctx) error {
	var req models.UserCreate
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	id, err := repository.CreateUser(req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "User creation failed"})
	}

	return c.Status(201).JSON(fiber.Map{
		"message":    "New user successfully created by admin",
		"created_id": id,
		"role":       req.Role,
		"clearance":  req.ClearanceLevel,
	})
}

func UpdateUserHandler(c *fiber.Ctx) error {
	userID := c.Params("id")

	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "missing userID in the url",
		})
	}

	var input models.UpdateUser
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid update payload",
		})
	}

	input.ID = userID

	err := repository.UpdateUser(input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "User updated successfully",
	})
}

func DeleteUserHandler(c *fiber.Ctx) error {
	userID := c.Params("id")

	err := repository.DeleteUser(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete user"})
	}

	return c.JSON(fiber.Map{"message": "User marked as deleted"})
}
