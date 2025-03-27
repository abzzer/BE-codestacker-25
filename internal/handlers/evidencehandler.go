package handlers

import (
	"strconv"
	"strings"

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

	evidenceID, err := repository.AddTextEvidence(req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not add text evidence"})
	}

	return c.JSON(fiber.Map{
		"message":     "Text evidence added",
		"evidence_id": evidenceID,
	})
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

	evidenceID, err := repository.AddTextEvidence(e)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not insert image evidence"})
	}

	return c.JSON(fiber.Map{
		"message":     "Image evidence added successfully",
		"evidence_id": evidenceID,
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

func UpdateEvidence(c *fiber.Ctx) error {
	evidenceIDStr := c.Params("evidenceid")
	evidenceID, err := strconv.Atoi(evidenceIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid evidence ID"})
	}

	currType, err := repository.GetEvidenceTypeByID(evidenceID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Evidence not found"})
	}

	if currType == models.EvidenceText {
		var payload struct {
			Content string `json:"content"`
		}
		if err := c.BodyParser(&payload); err != nil || strings.TrimSpace(payload.Content) == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid or empty text content"})
		}

		err = repository.UpdateEvidenceContent(evidenceID, payload.Content, "")
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to update text evidence"})
		}

		return c.JSON(fiber.Map{"message": "Text evidence updated successfully"})

	} else if currType == models.EvidenceImage {
		file, err := c.FormFile("image")
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Image file required for update"})
		}

		objectName, _, size, err := repository.UploadImageToMinio(file)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Image upload to MinIO failed"})
		}

		err = repository.UpdateEvidenceContent(evidenceID, objectName, size)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to update image evidence"})
		}

		return c.JSON(fiber.Map{"message": "Image evidence updated successfully"})
	}

	return c.Status(400).JSON(fiber.Map{"error": "Unsupported evidence type"})
}
