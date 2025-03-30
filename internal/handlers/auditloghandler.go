package handlers

import (
	"github.com/abzzer/BE-codestacker-25/internal/repository"
	"github.com/gofiber/fiber/v2"
)

func GetAuditLogs(c *fiber.Ctx) error {
	logs, err := repository.FetchAuditLogs()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch audit logs"})
	}

	return c.JSON(logs)
}
