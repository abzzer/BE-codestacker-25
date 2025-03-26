package handlers

import (
	"strconv"

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

	objectName, url, size, err := repository.UploadImageToMinio(file)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "MinIO upload failed"})
	}

	e := models.EvidenceTextRequest{
		CaseNumber: caseNumber,
		OfficerID:  officerID,
		Type:       models.EvidenceImage,
		Content:    objectName,
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

func GetEvidenceHandler(c *fiber.Ctx) error {
	evidenceIDParam := c.Params("evidenceid")
	evidenceID, err := strconv.Atoi(evidenceIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid evidence ID",
		})
	}

	evidence, err := repository.GetEvidenceByID(evidenceID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Evidence not found",
		})
	}

	if evidence.Type == "image" {
		return c.JSON(fiber.Map{
			"type":    "image",
			"remarks": evidence.Remarks,
			"size":    evidence.Size,
		})
	}

	return c.JSON(fiber.Map{
		"type":    "text",
		"remarks": evidence.Remarks,
		"content": evidence.Content,
	})
}

func GetImageEvidenceHandler(c *fiber.Ctx) error {
	evidenceIDStr := c.Params("evidenceid")
	evidenceID, err := strconv.Atoi(evidenceIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid evidence ID"})
	}

	reader, currType, err := repository.GetImageByID(evidenceID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	c.Set("Content-Type", currType)
	return c.SendStream(reader)
}
