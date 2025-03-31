package repository

import (
	"context"

	"github.com/abzzer/BE-codestacker-25/internal/database"
	"github.com/abzzer/BE-codestacker-25/internal/models"
)

func GetFullCaseDetails(caseNumber string) (*models.FullCaseDetails, error) {
	ctx := context.Background()

	base, err := GetCaseDetails(caseNumber)
	if err != nil {
		return nil, err
	}

	assigneeRows, err := database.DB.Query(ctx, `
		SELECT id, name, role, clearance_level FROM users
		WHERE id IN (
			SELECT user_id FROM case_assignees WHERE case_number = $1
		)
	`, caseNumber)
	if err != nil {
		return nil, err
	}
	defer assigneeRows.Close()

	var assignees []models.User
	for assigneeRows.Next() {
		var u models.User
		err := assigneeRows.Scan(&u.ID, &u.Name, &u.Role, &u.ClearanceLevel)
		if err == nil {
			assignees = append(assignees, u)
		}
	}

	evRows, err := database.DB.Query(ctx, `
		SELECT id, type, remarks, content, size
		FROM evidence
		WHERE case_number = $1 AND deleted = FALSE
	`, caseNumber)
	if err != nil {
		return nil, err
	}
	defer evRows.Close()

	var evidences []models.EvidenceWithID
	for evRows.Next() {
		var ev models.EvidenceWithID
		if err := evRows.Scan(&ev.ID, &ev.Type, &ev.Remarks, &ev.Content, &ev.Size); err == nil {
			evidences = append(evidences, ev)
		}
	}

	pRows, err := database.DB.Query(ctx, `
		SELECT id, type, name, age, gender, role
		FROM persons
		WHERE case_number = $1
	`, caseNumber)
	if err != nil {
		return nil, err
	}
	defer pRows.Close()

	var people []models.Person
	for pRows.Next() {
		var p models.Person
		if err := pRows.Scan(&p.ID, &p.Type, &p.Name, &p.Age, &p.Gender, &p.Role); err == nil {
			people = append(people, p)
		}
	}

	return &models.FullCaseDetails{
		CaseDetailsResponse: *base,
		Assignees:           assignees,
		Evidence:            evidences,
		People:              people,
	}, nil
}
