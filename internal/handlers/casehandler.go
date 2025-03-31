package handlers

import (
	"strings"

	"github.com/abzzer/BE-codestacker-25/internal/models"
	"github.com/abzzer/BE-codestacker-25/internal/repository"
	"github.com/gofiber/fiber/v2"
)

func CreateCaseHandler(c *fiber.Ctx) error {
	var req models.CaseRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "request body not valid",
		})
	}

	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Logged in users only -> You are not athorised",
		})
	}

	caseNumber, err := repository.CreateCase(req, userID.(string))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "couldn't create case"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":     "Case was created successfully here is your case number. No one is assigned this case yet.",
		"case_number": caseNumber,
	})
}

func UpdateCaseHandler(c *fiber.Ctx) error {
	caseID := c.Params("caseid")
	if caseID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing case ID in URL",
		})
	}

	var input models.CaseUpdate
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	input.CaseNumber = caseID

	if err := repository.UpdateCase(input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Case updated successfully",
	})
}

// FIND WAY TO STREAMLINE THIS METHOD SEEMS TO BE INEFFICIENT
func UpdateCaseStatusHandler(c *fiber.Ctx) error {
	caseID := c.Params("caseid")
	if caseID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing case ID in URL",
		})
	}

	var input models.CaseStatusUpdate
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON payload",
		})
	}

	// Check if we actually have status

	input.CaseNumber = caseID

	if err := repository.UpdateCaseStatus(input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Case status updated successfully",
	})
}

func AddPersonHandler(c *fiber.Ctx) error {
	caseID := c.Params("caseid")
	if caseID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing case ID in URL",
		})
	}

	var person models.PersonRequest
	if err := c.BodyParser(&person); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON payload",
		})
	}

	person.CaseNumber = caseID

	personID, err := repository.AddPersonToCase(person)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add person to case",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":   "Person added successfully",
		"person_id": personID,
	})
}

func GetPartialCaseDetailsHandler(c *fiber.Ctx) error {
	caseID := c.Params("caseid")
	if caseID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Case ID missing",
		})
	}

	caseID = strings.ToUpper(caseID)

	caseLevel, err := repository.GetCaseLevelByNumber(caseID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Case not found",
		})
	}

	userRole := c.Locals("role").(string)
	userClearance := c.Locals("clearance").(string)

	if err := repository.CheckClearance(userRole, userClearance, caseLevel); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	details, err := repository.GetCaseDetails(caseID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch case details",
		})
	}

	return c.JSON(details)
}

func GetFullCaseDetailsHandler(c *fiber.Ctx) error {
	caseID := c.Params("caseid")
	if caseID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Case ID missing",
		})
	}
	caseID = strings.ToUpper(caseID)

	caseLevel, err := repository.GetCaseLevelByNumber(caseID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Case not found",
		})
	}

	userRole := c.Locals("role").(string)
	userClearance := c.Locals("clearance").(string)

	if err := repository.CheckClearance(userRole, userClearance, caseLevel); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	fullDetails, err := repository.GetFullCaseDetails(caseID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch full case details",
		})
	}

	return c.JSON(fullDetails)
}

func AddOfficerToCaseHandler(c *fiber.Ctx) error {
	caseID := c.Params("caseid")
	if caseID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing case ID in URL",
		})
	}

	var req struct {
		UserID string `json:"user_id"`
	}

	if err := c.BodyParser(&req); err != nil || req.UserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required in the JSON body",
		})
	}

	if err := repository.AssignUserToCase(req.UserID, caseID); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "User successfully assigned to case",
	})
}
