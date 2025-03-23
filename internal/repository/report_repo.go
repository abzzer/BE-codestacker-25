package repository

import (
	"context"

	"github.com/abzzer/BE-codestacker-25/internal/database"
	"github.com/abzzer/BE-codestacker-25/internal/models"
)

func SubmitCrimeReport(report models.CrimeReportRequest) (int, error) {
	query := `
		INSERT INTO reports (email, civil_id, name, description, area, city)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING report_id;
	`
	var reportID int
	err := database.DB.QueryRow(context.Background(), query, report.Email, report.CivilID, report.Name, report.Description, report.Area, report.City).Scan(&reportID)
	if err != nil {
		return 0, err
	}

	return reportID, nil
}
