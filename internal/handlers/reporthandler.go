package handlers

import (
	"strconv"

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

func GetAllReports(c *fiber.Ctx) error {
	reports, err := repository.GetAllReports()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch reports"})
	}
	return c.JSON(reports)
}

func LinkReportToCase(c *fiber.Ctx) error {
	reportID, err := strconv.Atoi(c.Params("reportID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid report ID"})
	}

	var request struct {
		CaseNumber string `json:"case_number"`
	}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	err = repository.LinkReportToCase(reportID, request.CaseNumber)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to link report to case"})
	}

	return c.JSON(fiber.Map{"message": "Report linked to case successfully"})
}

func CheckReportStatus(c *fiber.Ctx) error {
	reportID, err := strconv.Atoi(c.Params("reportID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid report ID"})
	}

	status, err := repository.GetReportStatus(reportID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": status})
}
