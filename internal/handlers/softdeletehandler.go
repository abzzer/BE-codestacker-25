package handlers

import (
	"strconv"

	"github.com/abzzer/BE-codestacker-25/internal/repository"
	"github.com/gofiber/fiber/v2"
)

func DeleteEvidence(c *fiber.Ctx) error {
	evidenceIDStr := c.Params("evidenceid")
	evidenceID, err := strconv.Atoi(evidenceIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid evidence ID"})
	}

	userID := c.Locals("user_id").(string)

	err = repository.SoftDeleteEvidence(evidenceID, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Evidence soft-deleted successfully and audit log added",
	})
}
