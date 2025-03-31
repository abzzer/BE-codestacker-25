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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid evidence ID"})
	}

	userID := c.Locals("user_id").(string)

	switch c.Method() {

	case fiber.MethodPost:
		state.SetStatus(evidenceID, state.StatusInitiated)
		return c.JSON(fiber.Map{
			"message": "Are you sure you want to permanently delete Evidence ID: " + evidenceIDStr + "? Send PATCH with confirm: 'yes' in json.",
		})

	case fiber.MethodPatch:
		var confirm struct {
			Confirm string `json:"confirm"`
		}
		if err := c.BodyParser(&confirm); err != nil || confirm.Confirm != "yes" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Confirmation failed. Send PATCH with JSON with key confirm and value yes",
			})
		}
		state.SetStatus(evidenceID, state.StatusConfirmed)
		return c.JSON(fiber.Map{
			"message": "Confirmation accepted. Now send DELETE to complete hard deletion.",
		})

	case fiber.MethodDelete:
		if state.GetStatus(evidenceID) != state.StatusConfirmed {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Confirmation required. First send PATCH with JSON key 'confirm' and value yes before DELETE.",
			})
		}

		state.SetStatus(evidenceID, state.StatusDeleting)

		if err := repository.HardDeleteEvidence(evidenceID, userID); err != nil {
			state.SetStatus(evidenceID, state.StatusFailed)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		state.SetStatus(evidenceID, state.StatusDone)
		return c.JSON(fiber.Map{
			"message":     "Evidence hard-deleted successfully.",
			"evidence_id": evidenceID,
			"status":      state.StatusDone,
		})

	default:
		return c.Status(fiber.StatusMethodNotAllowed).JSON(fiber.Map{
			"error": "Unsupported method. Use POST to start, PATCH to confirm, DELETE to execute.",
		})
	}
}
