package handlers

import (
	"github.com/abzzer/BE-codestacker-25/internal/models"
	"github.com/abzzer/BE-codestacker-25/internal/repository"
	"github.com/gofiber/fiber/v2"
)

func SubmitCrimeReportHandler(c *fiber.Ctx) error {
	var req models.CrimeReportRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	reportID, err := repository.SubmitCrimeReport(req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to submit report"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":   "Report submitted successfully. Please keep your report ID to check status.",
		"report_id": reportID,
	})
}
