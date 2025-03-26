package repository

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/abzzer/BE-codestacker-25/internal/database"
	"github.com/abzzer/BE-codestacker-25/internal/models"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

func AddTextEvidence(e models.EvidenceTextRequest) error {
	query := `
		INSERT INTO evidence (case_number, officer_id, type, content, size, remarks)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := database.DB.Exec(context.Background(), query,
		e.CaseNumber, e.OfficerID, e.Type, e.Content, e.Size, e.Remarks,
	)
	return err
}

func UploadImageToMinio(file *multipart.FileHeader) (string, string, error) {
	bucket := "evidence-bucket"
	fileExt := filepath.Ext(file.Filename)
	fileID := uuid.New().String()
	objectName := fmt.Sprintf("evidence/%s%s", fileID, fileExt)

	src, err := file.Open()
	if err != nil {
		return "", "", err
	}
	defer src.Close()

	contentType := file.Header.Get("Content-Type")
	fileSize := file.Size

	_, err = database.MinioClient.PutObject(context.Background(), bucket, objectName, src, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", "", err
	}

	url := fmt.Sprintf("http://%s/%s/%s",
		os.Getenv("MINIO_ENDPOINT"), bucket, objectName)

	return url, fmt.Sprintf("%d bytes", fileSize), nil
}

func GetEvidenceByID(id int) (*models.EvidenceFromID, error) {
	query := `
		SELECT type, remarks, content, size
		FROM evidence
		WHERE id = $1 AND deleted = FALSE;
	`

	var ev models.EvidenceFromID
	err := database.DB.QueryRow(context.Background(), query, id).Scan(&ev.Type, &ev.Remarks, &ev.Content, &ev.Size)
	if err != nil {
		return nil, errors.New("evidence not found")
	}

	return &ev, nil
}
