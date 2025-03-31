package repository

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"

	"github.com/abzzer/BE-codestacker-25/internal/database"
	"github.com/jung-kurt/gofpdf"
)

func GenerateCasePDF(caseNumber string) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)

	details, err := GetCaseDetails(caseNumber)
	if err != nil {
		return nil, err
	}

	pdf.Cell(40, 10, fmt.Sprintf("Case Report - %s", caseNumber))
	pdf.Ln(12)
	pdf.SetFont("Arial", "", 12)
	createdAtStr := details.CreatedAt.Format("2006-01-02 15:04:05")

	pdf.MultiCell(0, 8, fmt.Sprintf(`Case Name: %s
Description: %s
Location: %s, %s
Created By: %s
Created At: %s
Case Type: %s
Level: %s
Status: %s
Reported By: %d
Assignees: %d
Evidences: %d
Suspects: %d
Victims: %d
Witnesses: %d
`, details.CaseName, details.Description, details.Area, details.City, details.CreatedBy, createdAtStr, details.CaseType, details.Level, details.Status,
		details.ReportedBy, details.NumAssignees, details.NumEvidences, details.NumSuspects, details.NumVictims, details.NumWitnesses), "", "", false)

	pdf.Ln(4)
	pdf.SetFont("Arial", "B", 13)
	pdf.Cell(0, 10, "Case Assignees")
	pdf.Ln(10)

	rows, err := database.DB.Query(context.Background(), `
		SELECT u.id, u.name, u.role, u.clearance_level
		FROM case_assignees ca
		JOIN users u ON u.id = ca.user_id
		WHERE ca.case_number = $1
	`, caseNumber)
	if err == nil {
		pdf.SetFont("Arial", "", 11)
		for rows.Next() {
			var id, name, role, clearance string
			_ = rows.Scan(&id, &name, &role, &clearance)
			pdf.MultiCell(0, 6, fmt.Sprintf("ID: %s | Name: %s | Role: %s | Clearance: %s", id, name, role, clearance), "", "", false)
		}
		rows.Close()
	}

	pdf.Ln(6)
	pdf.SetFont("Arial", "B", 13)
	pdf.Cell(0, 10, "People Involved")
	pdf.Ln(10)

	rows, err = database.DB.Query(context.Background(), `
		SELECT type, name, age, gender, role
		FROM persons
		WHERE case_number = $1
	`, caseNumber)
	if err == nil {
		pdf.SetFont("Arial", "", 11)
		for rows.Next() {
			var ptype, name, gender, role string
			var age int
			_ = rows.Scan(&ptype, &name, &age, &gender, &role)
			pdf.MultiCell(0, 6, fmt.Sprintf("Type: %s | Name: %s | Age: %d | Gender: %s | Role: %s", ptype, name, age, gender, role), "", "", false)
		}
		rows.Close()
	}

	pdf.Ln(6)
	pdf.SetFont("Arial", "B", 13)
	pdf.Cell(0, 10, "Citizen Reports Linked")
	pdf.Ln(10)

	rows, err = database.DB.Query(context.Background(), `
		SELECT name, email, civil_id, description
		FROM reports
		WHERE case_number = $1
	`, caseNumber)
	if err == nil {
		pdf.SetFont("Arial", "", 11)
		for rows.Next() {
			var name, email, civilID, desc string
			_ = rows.Scan(&name, &email, &civilID, &desc)
			pdf.MultiCell(0, 6, fmt.Sprintf("Name: %s | Email: %s | Civil ID: %s\nDescription: %s", name, email, civilID, desc), "", "", false)
			pdf.Ln(2)
		}
		rows.Close()
	}

	pdf.Ln(6)
	pdf.SetFont("Arial", "B", 13)
	pdf.Cell(0, 10, "Text Evidence")
	pdf.Ln(10)

	rows, err = database.DB.Query(context.Background(), `
		SELECT content, remarks
		FROM evidence
		WHERE case_number = $1 AND type = 'text' AND deleted = FALSE
	`, caseNumber)
	if err == nil {
		pdf.SetFont("Arial", "", 11)
		for rows.Next() {
			var content, remarks string
			_ = rows.Scan(&content, &remarks)
			pdf.MultiCell(0, 6, fmt.Sprintf("Remarks: %s\nContent: %s", remarks, content), "", "", false)
			pdf.Ln(2)
		}
		rows.Close()
	}

	pdf.Ln(6)
	pdf.SetFont("Arial", "B", 13)
	pdf.Cell(0, 10, "Image Evidence")
	pdf.Ln(10)

	rows, err = database.DB.Query(context.Background(), `
		SELECT id, content, remarks, size
		FROM evidence
		WHERE case_number = $1 AND type = 'image' AND deleted = FALSE
	`, caseNumber)

	if err == nil {
		pdf.SetFont("Arial", "", 11)
		for rows.Next() {
			var eid int
			var content, remarks, size string
			_ = rows.Scan(&eid, &content, &remarks, &size)
			pdf.MultiCell(0, 6, fmt.Sprintf("ID: %d | Remarks: %s\nImage Path: %s\nSize: %s", eid, remarks, content, size), "", "", false)
			pdf.Ln(2)

			reader, _, err := GetImageByID(eid)
			if err == nil {
				img, _, err := image.Decode(reader)
				if err == nil {
					tmp := new(bytes.Buffer)
					jpeg.Encode(tmp, img, nil)
					imgID := fmt.Sprintf("img%d", eid)
					pdf.RegisterImageOptionsReader(imgID, gofpdf.ImageOptions{ImageType: "JPEG"}, tmp)
					pdf.ImageOptions(imgID, 10, pdf.GetY(), 60, 0, false, gofpdf.ImageOptions{ImageType: "JPEG"}, 0, "")
					pdf.Ln(5)
				}
			}
		}
		rows.Close()
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
