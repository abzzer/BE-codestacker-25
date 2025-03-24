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
