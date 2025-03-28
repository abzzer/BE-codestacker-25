package handlers

import (
	"strconv"

	"github.com/abzzer/BE-codestacker-25/internal/repository"
	"github.com/abzzer/BE-codestacker-25/internal/state"
	"github.com/gofiber/fiber/v2"
)

func SoftDeleteEvidence(c *fiber.Ctx) error {
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
func HardDeleteEvidence(c *fiber.Ctx) error {
	evidenceIDStr := c.Params("evidenceid")
	evidenceID, err := strconv.Atoi(evidenceIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid evidence ID"})
	}

	userID := c.Locals("user_id").(string)

	// Check request method
	switch c.Method() {
	case fiber.MethodPost:
		// Step 1: Prompt for confirmation
		state.SetStatus(userID, evidenceID, state.StatusInitiated)
		return c.JSON(fiber.Map{
			"message": "Are you sure you want to permanently delete Evidence ID: " + evidenceIDStr + "? Send POST again with {\"confirm\":\"yes\"} or send DELETE to proceed.",
		})

	case fiber.MethodDelete:
		// Step 2: Only allow if previously confirmed
		status := state.GetStatus(userID, evidenceID)
		if status != state.StatusConfirmed {
			return c.Status(400).JSON(fiber.Map{
				"error": "Confirmation required. First send POST with {\"confirm\":\"yes\"} before DELETE.",
			})
		}

		state.SetStatus(userID, evidenceID, state.StatusDeleting)

		err := repository.HardDeleteEvidence(evidenceID, userID)
		if err != nil {
			state.SetStatus(userID, evidenceID, state.StatusFailed)
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		state.SetStatus(userID, evidenceID, state.StatusDone)
		return c.JSON(fiber.Map{
			"message":     "Evidence hard-deleted successfully.",
			"evidence_id": evidenceID,
			"status":      state.StatusDone,
		})

	case fiber.MethodPatch:
		// Step 1b: User confirms deletion
		var confirm struct {
			Confirm string `json:"confirm"`
		}

		if err := c.BodyParser(&confirm); err != nil || confirm.Confirm != "yes" {
			return c.Status(400).JSON(fiber.Map{
				"error": "Confirmation failed. Send PATCH with JSON {\"confirm\": \"yes\"}",
			})
		}

		state.SetStatus(userID, evidenceID, state.StatusConfirmed)
		return c.JSON(fiber.Map{
			"message": "Confirmation accepted. Now send DELETE to complete hard deletion.",
		})

	default:
		return c.Status(405).JSON(fiber.Map{
			"error": "Unsupported method. Use POST to start, PATCH to confirm, DELETE to execute.",
		})
	}
}
