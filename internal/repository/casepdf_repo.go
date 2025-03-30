package repository

import (
	"bytes"
	"context"
	"fmt"

	"github.com/abzzer/BE-codestacker-25/internal/database"
	"github.com/jung-kurt/gofpdf"
)

func GenerateCasePDF(caseNumber string) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)

	var caseName, description, area, city, createdBy, caseType, level, status string
	err := database.DB.QueryRow(context.Background(), `
		SELECT case_name, description, area, city, created_by, case_type, level, status
		FROM cases
		WHERE case_number = $1
	`, caseNumber).Scan(&caseName, &description, &area, &city, &createdBy, &caseType, &level, &status)

	if err != nil {
		return nil, err
	}

	pdf.Cell(40, 10, fmt.Sprintf("Case Report - %s", caseNumber))
	pdf.Ln(12)
	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(0, 8, fmt.Sprintf(`
Case Name: %s
Description: %s
Location: %s, %s
Created By: %s
Case Type: %s
Level: %s
Status: %s
`, caseName, description, area, city, createdBy, caseType, level, status), "", "", false)

	// Fetch Persons
	pdf.Ln(5)
	pdf.SetFont("Arial", "B", 13)
	pdf.Cell(0, 10, "People Involved")
	pdf.Ln(10)
	rows, _ := database.DB.Query(context.Background(), `
		SELECT type, name, age, gender, role
		FROM persons
		WHERE case_number = $1
	`, caseNumber)
	defer rows.Close()

	pdf.SetFont("Arial", "", 11)
	for rows.Next() {
		var ptype, name, gender, role string
		var age int
		_ = rows.Scan(&ptype, &name, &age, &gender, &role)
		pdf.MultiCell(0, 6, fmt.Sprintf("Type: %s | Name: %s | Age: %d | Gender: %s | Role: %s", ptype, name, age, gender, role), "", "", false)
	}

	// Fetch Text Evidence
	pdf.Ln(5)
	pdf.SetFont("Arial", "B", 13)
	pdf.Cell(0, 10, "Evidence (Text)")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 11)
	rows2, _ := database.DB.Query(context.Background(), `
		SELECT content, remarks
		FROM evidence
		WHERE case_number = $1 AND type = 'text' AND deleted = FALSE
	`, caseNumber)
	defer rows2.Close()

	for rows2.Next() {
		var content, remarks string
		_ = rows2.Scan(&content, &remarks)
		pdf.MultiCell(0, 6, fmt.Sprintf("Remarks: %s\nContent: %s", remarks, content), "", "", false)
		pdf.Ln(2)
	}

	// Finish
	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
