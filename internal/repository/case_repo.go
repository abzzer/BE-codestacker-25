package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/abzzer/BE-codestacker-25/internal/database"
	"github.com/abzzer/BE-codestacker-25/internal/models"
)

func truncateDescription(desc string) string {
	if len(desc) <= 100 {
		return desc
	}
	limit := 96 //because "_..." <- 4 chars
	lastSpace := strings.LastIndex(desc[:limit], " ")
	if lastSpace == -1 {
		return " ..." // if we don't have a space basically
	}
	return desc[:lastSpace] + " ..."
}

func CreateCase(caseReq models.CaseRequest, createdBy string) (string, error) {
	desc := truncateDescription(caseReq.Description)

	query := `
		INSERT INTO cases (case_name, description, area, city, created_by, level)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING case_number;
	`

	var caseNumber string
	err := database.DB.QueryRow(context.Background(), query, caseReq.CaseName, desc, caseReq.Area, caseReq.City, createdBy, caseReq.Level).Scan(&caseNumber)

	if err != nil {
		return "", err
	}
	return caseNumber, nil
}

func UpdateCase(input models.CaseUpdate) error {
	if input.CaseNumber == "" {
		return errors.New("case_number is required")
	}

	setClauses := []string{}
	args := []interface{}{}
	argIdx := 1

	if input.CaseName != nil {
		setClauses = append(setClauses, fmt.Sprintf("case_name = $%d", argIdx))
		args = append(args, *input.CaseName)
		argIdx++
	}

	if input.Description != nil {
		truncated := truncateDescription(*input.Description)
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", argIdx))
		args = append(args, truncated)
		argIdx++
	}

	if input.Area != nil {
		setClauses = append(setClauses, fmt.Sprintf("area = $%d", argIdx))
		args = append(args, *input.Area)
		argIdx++
	}

	if input.City != nil {
		setClauses = append(setClauses, fmt.Sprintf("city = $%d", argIdx))
		args = append(args, *input.City)
		argIdx++
	}

	if input.Level != nil {
		setClauses = append(setClauses, fmt.Sprintf("level = $%d", argIdx))
		args = append(args, *input.Level)
		argIdx++
	}

	if input.Status != nil {
		setClauses = append(setClauses, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, *input.Status)
		argIdx++
	}

	if len(setClauses) == 0 {
		return errors.New("no valid fields provided to update")
	}

	args = append(args, input.CaseNumber)

	query := fmt.Sprintf(`
		UPDATE cases SET %s
		WHERE case_number = $%d;
	`, strings.Join(setClauses, ", "), argIdx)

	_, err := database.DB.Exec(context.Background(), query, args...)
	return err
}

func UpdateCaseStatus(input models.CaseStatusUpdate) error {
	if input.CaseNumber == "" {
		return errors.New("case_number is required")
	}

	query := `UPDATE cases SET status = $1 WHERE case_number = $2;`
	_, err := database.DB.Exec(context.Background(), query, input.Status, input.CaseNumber)
	return err
}

func AddPersonToCase(person models.PersonRequest) (int, error) {
	query := `
		INSERT INTO persons (case_number, type, name, age, gender, role)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id;
	`

	var personID int
	err := database.DB.QueryRow(context.Background(), query, person.CaseNumber, person.Type, person.Name, person.Age, person.Gender, person.Role).Scan(&personID)

	if err != nil {
		return 0, err
	}

	return personID, nil
}

func GetCaseLevelByNumber(caseNumber string) (models.CaseLevel, error) {
	query := `SELECT level FROM cases WHERE case_number = $1`
	var level string

	err := database.DB.QueryRow(context.Background(), query, caseNumber).Scan(&level)
	if err != nil {
		return "", errors.New("could not find case or retrieve level")
	}

	return models.CaseLevel(level), nil
}

func GetCaseDetails(caseNumber string) (*models.CaseDetailsResponse, error) {
	query := `
	SELECT 
		c.case_number, c.case_name, c.description, c.area, c.city,
		c.created_by, c.created_at, c.case_type, c.level, c.status,
		(SELECT COUNT(*) FROM reports WHERE case_number = c.case_number) AS reported_by,
		(SELECT COUNT(*) FROM case_assignees WHERE case_number = c.case_number) AS num_assignees,
		(SELECT COUNT(*) FROM evidence WHERE case_number = c.case_number AND deleted = FALSE) AS num_evidences,
		(SELECT COUNT(*) FROM persons WHERE case_number = c.case_number AND type = 'suspect') AS num_suspects,
		(SELECT COUNT(*) FROM persons WHERE case_number = c.case_number AND type = 'victim') AS num_victims,
		(SELECT COUNT(*) FROM persons WHERE case_number = c.case_number AND type = 'witness') AS num_witnesses
	FROM cases c
	WHERE c.case_number = $1;
	`

	var result models.CaseDetailsResponse
	var levelStr, statusStr string

	err := database.DB.QueryRow(context.Background(), query, caseNumber).Scan(&result.CaseNumber,
		&result.CaseName, &result.Description, &result.Area, &result.City, &result.CreatedBy,
		&result.CreatedAt, &result.CaseType, &levelStr, &statusStr, &result.ReportedBy, &result.NumAssignees,
		&result.NumEvidences, &result.NumSuspects, &result.NumVictims, &result.NumWitnesses)

	if err != nil {
		return nil, err
	}

	result.Level = models.CaseLevel(levelStr)
	result.Status = models.CaseStatus(statusStr)

	return &result, nil
}

func AssignUserToCase(targetUserID, caseNumber string) error {
	ctx := context.Background()

	var userRole string
	var userClearance string
	var caseClearance string

	err := database.DB.QueryRow(ctx, `
		SELECT role, clearance_level FROM users
		WHERE id = $1 AND deleted = FALSE
	`, targetUserID).Scan(&userRole, &userClearance)
	if err != nil {
		return errors.New("target user not found")
	}

	err = database.DB.QueryRow(ctx, `
		SELECT level FROM cases
		WHERE case_number = $1
	`, caseNumber).Scan(&caseClearance)
	if err != nil {
		return errors.New("case not found")
	}

	if userRole == "officer" {
		if !IsClearanceSufficient(models.CaseLevel(userClearance), models.CaseLevel(caseClearance)) {
			return errors.New("officer's clearance level is insufficient for this case")
		}
	}

	_, err = database.DB.Exec(ctx, `
		INSERT INTO case_assignees (case_number, user_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`, caseNumber, targetUserID)
	if err != nil {
		return errors.New("failed to assign user to case")
	}

	return nil
}
