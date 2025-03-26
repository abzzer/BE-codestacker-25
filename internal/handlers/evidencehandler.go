package handlers

import (
	"github.com/abzzer/BE-codestacker-25/internal/models"
	"github.com/abzzer/BE-codestacker-25/internal/repository"
	"github.com/gofiber/fiber/v2"
)

func AddTextEvidenceHandler(c *fiber.Ctx) error {
	officerID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User ID not found in token"})
	}

	var req models.EvidenceTextRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if req.Type != models.EvidenceText {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid type for this endpoint"})
	}

	req.OfficerID = officerID

	if err := repository.AddTextEvidence(req); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not add text evidence"})
	}

	return c.JSON(fiber.Map{"message": "Text evidence added"})
}

func AddImageEvidenceHandler(c *fiber.Ctx) error {
	officerID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User ID not found in token"})
	}

	caseNumber := c.FormValue("case_number")
	remarks := c.FormValue("remarks")

	if caseNumber == "" {
		return c.Status(400).JSON(fiber.Map{"error": "case_number is required"})
	}

	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Image file is required"})
	}

	url, size, err := repository.UploadImageToMinio(file)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "MinIO upload failed"})
	}

	e := models.EvidenceTextRequest{
		CaseNumber: caseNumber,
		OfficerID:  officerID,
		Type:       models.EvidenceImage,
		Content:    url,
		Size:       size,
		Remarks:    remarks,
	}

	if err := repository.AddTextEvidence(e); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not insert image evidence"})
	}

	return c.JSON(fiber.Map{
		"message":     "Image evidence added successfully",
		"minio_url":   url,
		"contentSize": size,
	})
}
