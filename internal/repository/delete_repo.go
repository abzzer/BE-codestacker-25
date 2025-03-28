package repository

import (
	"context"
	"errors"

	"github.com/abzzer/BE-codestacker-25/internal/database"
	"github.com/abzzer/BE-codestacker-25/internal/models"
	"github.com/minio/minio-go/v7"
)

func UpdateEvidenceContent(id int, content, size string) error {
	query := `UPDATE evidence SET content = $1, size = $2 WHERE id = $3 AND deleted = FALSE;`
	_, err := database.DB.Exec(context.Background(), query, content, size, id)
	return err
}

func SoftDeleteEvidence(evidenceID int, userID string) error {
	query := `
		UPDATE evidence
		SET deleted = TRUE
		WHERE id = $1 AND deleted = FALSE
		RETURNING id;
	`

	var id int
	err := database.DB.QueryRow(context.Background(), query, evidenceID).Scan(&id)
	if err != nil {
		return errors.New("failed to soft delete or evidence already deleted")
	}

	auditQuery := `
		INSERT INTO audit_logs (action, evidence_id, user_id)
		VALUES ($1, $2, $3)
	`

	_, err = database.DB.Exec(context.Background(), auditQuery, models.AuditSoftDeleted, id, userID)
	if err != nil {
		return errors.New("soft deleted but failed to log audit")
	}

	return nil
}

func HardDeleteEvidence(evidenceID int, userID string) error {
	var evType models.EvidenceType
	var content string

	query := `
		SELECT type, content FROM evidence
		WHERE id = $1
	`
	err := database.DB.QueryRow(context.Background(), query, evidenceID).Scan(&evType, &content)
	if err != nil {
		return errors.New("evidence not found")
	}

	if evType == models.EvidenceImage {
		err := database.MinioClient.RemoveObject(context.Background(), "evidence-bucket", content, minio.RemoveObjectOptions{})
		if err != nil {
			return errors.New("failed to delete image from MinIO")
		}
	}

	deleteQuery := `
		DELETE FROM evidence
		WHERE id = $1
	`
	_, err = database.DB.Exec(context.Background(), deleteQuery, evidenceID)
	if err != nil {
		return errors.New("failed to delete evidence from database")
	}

	auditQuery := `
		INSERT INTO audit_logs (action, evidence_id, user_id)
		VALUES ($1, $2, $3)
	`
	_, err = database.DB.Exec(context.Background(), auditQuery, models.AuditHardDeleted, evidenceID, userID)
	if err != nil {
		return errors.New("evidence deleted but failed to write audit log")
	}

	return nil
}
