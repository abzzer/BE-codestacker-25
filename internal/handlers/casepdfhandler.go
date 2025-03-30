package handlers

import (
	"github.com/abzzer/BE-codestacker-25/internal/repository"
	"github.com/gofiber/fiber/v2"
)

func GenerateCasePDFHandler(c *fiber.Ctx) error {
	caseID := c.Params("caseid")
	if caseID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing case ID"})
	}

	pdfBytes, err := repository.GenerateCasePDF(caseID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate PDF: " + err.Error()})
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "attachment; filename=\"case_report_"+caseID+".pdf\"")
	return c.Send(pdfBytes)
}
