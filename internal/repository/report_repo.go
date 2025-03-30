package repository

import (
	"context"
	"errors"

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
func GetAllReports() ([]models.Report, error) {
	query := `SELECT report_id, email, civil_id, name, role, case_number, description, area, city FROM reports`
	rows, err := database.DB.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []models.Report
	for rows.Next() {
		var r models.Report
		if err := rows.Scan(&r.ReportID, &r.Email, &r.CivilID, &r.Name, &r.Role, &r.CaseNumber, &r.Description, &r.Area, &r.City); err != nil {
			return nil, err
		}
		reports = append(reports, r)
	}
	return reports, nil
}

func LinkReportToCase(reportID int, caseNumber string) error {
	query := `UPDATE reports SET case_number = $1 WHERE report_id = $2`
	_, err := database.DB.Exec(context.Background(), query, caseNumber, reportID)
	return err
}

func GetReportStatus(reportID int) (string, error) {
	var caseNumber *string
	query := `SELECT case_number FROM reports WHERE report_id = $1`
	err := database.DB.QueryRow(context.Background(), query, reportID).Scan(&caseNumber)
	if err != nil {
		return "", errors.New("report not found")
	}

	if caseNumber == nil {
		return "pending", nil
	}

	var status string
	caseQuery := `SELECT status FROM cases WHERE case_number = $1`
	err = database.DB.QueryRow(context.Background(), caseQuery, *caseNumber).Scan(&status)
	if err != nil {
		return "", errors.New("case related to report not found")
	}

	return status, nil
}
