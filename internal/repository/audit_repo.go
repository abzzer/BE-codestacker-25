package repository

import (
	"context"

	"github.com/abzzer/BE-codestacker-25/internal/database"
	"github.com/abzzer/BE-codestacker-25/internal/models"
)

func FetchAuditLogs() ([]models.AuditLog, error) {
	query := `
		SELECT id, action, evidence_id, user_id, timestamp
		FROM audit_logs
		ORDER BY timestamp DESC
	`

	rows, err := database.DB.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.AuditLog
	for rows.Next() {
		var log models.AuditLog
		err := rows.Scan(&log.ID, &log.Action, &log.EvidenceID, &log.UserID, &log.Timestamp)
		if err != nil {
			continue
		}
		logs = append(logs, log)
	}

	return logs, nil
}
