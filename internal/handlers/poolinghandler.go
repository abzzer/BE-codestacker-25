package handlers

import (
	"strconv"
	"time"

	"github.com/abzzer/BE-codestacker-25/internal/state"
	"github.com/gofiber/fiber/v2"
)

func LongPollDeleteStatus(c *fiber.Ctx) error {
	evidenceIDStr := c.Params("evidenceid")
	evidenceID, err := strconv.Atoi(evidenceIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid evidence ID"})
	}

	userID := c.Locals("user_id").(string)

	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return c.JSON(fiber.Map{
				"status":  state.GetStatus(userID, evidenceID),
				"message": "Timeout reached, no final status yet",
			})
		case <-ticker.C:
			current := state.GetStatus(userID, evidenceID)
			if current == state.StatusDone || current == state.StatusFailed {
				state.ClearStatus(userID, evidenceID)

				return c.JSON(fiber.Map{
					"status":  current,
					"message": "Deletion status resolved",
				})
			}
		}
	}
}
