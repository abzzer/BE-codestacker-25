package repository

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/abzzer/BE-codestacker-25/internal/database"
	"github.com/abzzer/BE-codestacker-25/internal/models"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

func AddTextEvidence(e models.EvidenceTextRequest) (int, error) {
	query := `
		INSERT INTO evidence (case_number, officer_id, type, content, size, remarks)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id;
	`
	var newID int
	err := database.DB.QueryRow(context.Background(), query, e.CaseNumber, e.OfficerID, e.Type, e.Content, e.Size, e.Remarks).Scan(&newID)
	return newID, err
}

func UploadImageToMinio(file *multipart.FileHeader) (string, string, string, error) {
	fileExt := filepath.Ext(file.Filename)
	fileID := uuid.New().String()
	objectName := fmt.Sprintf("evidence/%s%s", fileID, fileExt)

	src, err := file.Open()
	if err != nil {
		return "", "", "", err
	}
	defer src.Close()

	contentType := file.Header.Get("Content-Type")
	fileSize := file.Size

	_, err = database.MinioClient.PutObject(context.Background(), "evidence-bucket", objectName, src, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", "", "", err
	}

	url := fmt.Sprintf("http://%s/%s/%s",
		os.Getenv("MINIO_ENDPOINT"), "evidence-bucket", objectName)

	return objectName, url, fmt.Sprintf("%d bytes", fileSize), nil
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

func GetImageByID(id int) (io.Reader, string, error) {
	query := `
		SELECT type, content
		FROM evidence
		WHERE id = $1 AND deleted = FALSE;
	`

	var ev models.ImageFromID
	err := database.DB.QueryRow(context.Background(), query, id).Scan(&ev.Type, &ev.Content)
	if err != nil {
		return nil, "", errors.New("evidence not found")
	}

	if ev.Type != models.EvidenceImage {
		return nil, "", errors.New("evidence is not an image")
	}
	object, err := database.MinioClient.GetObject(
		context.Background(),
		"evidence-bucket",
		ev.Content,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return nil, "", errors.New("failed to access image in MinIO")
	}

	info, err := object.Stat()
	if err != nil {
		return nil, "", errors.New("could not retrieve object info")
	}

	if !strings.HasPrefix(info.ContentType, "image/") {
		return nil, "", errors.New("file in MinIO is not an image")
	}

	return object, info.ContentType, nil
}

func GetEvidenceTypeByID(id int) (models.EvidenceType, error) {
	var evType models.EvidenceType
	query := `SELECT type FROM evidence WHERE id = $1 AND deleted = FALSE;`
	err := database.DB.QueryRow(context.Background(), query, id).Scan(&evType)
	if err != nil {
		return "", errors.New("evidence not found")
	}
	return evType, nil
}

func GetTopTextEvidenceWords() ([]string, error) {
	query := `SELECT content FROM evidence WHERE type = 'text' AND deleted = FALSE`
	rows, err := database.DB.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var texts []string
	for rows.Next() {
		var content string
		if err := rows.Scan(&content); err == nil {
			texts = append(texts, content)
		}
	}

	stopWords := map[string]struct{}{
		"the": {}, "and": {}, "to": {}, "a": {}, "of": {}, "in": {}, "on": {}, "with": {},
		"at": {}, "by": {}, "for": {}, "an": {}, "was": {}, "is": {}, "were": {}, "had": {},
		"be": {}, "it": {}, "that": {}, "this": {}, "as": {}, "from": {}, "but": {}, "or": {},
		"are": {}, "before": {}, "after": {}, "same": {},
	}

	wordCount := make(map[string]int)
	for _, text := range texts {
		words := regexp.MustCompile(`\b[a-zA-Z]+\b`).FindAllString(strings.ToLower(text), -1)
		for _, word := range words {
			if _, skip := stopWords[word]; !skip {
				wordCount[word]++
			}
		}
	}

	type wordFreq struct {
		Word  string
		Count int
	}
	var sorted []wordFreq
	for w, c := range wordCount {
		sorted = append(sorted, wordFreq{w, c})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Count > sorted[j].Count
	})

	top := []string{}
	for i := 0; i < len(sorted) && i < 10; i++ {
		top = append(top, sorted[i].Word)
	}
	return top, nil
}

func ExtractURLsFromCase(caseNumber string) ([]string, error) {
	query := `
		SELECT content FROM evidence
		WHERE case_number = $1 AND type = 'text' AND deleted = FALSE
	`

	rows, err := database.DB.Query(context.Background(), query, caseNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []string
	regex := regexp.MustCompile(`https?://[^\s]+`)

	for rows.Next() {
		var content string
		if err := rows.Scan(&content); err != nil {
			continue
		}

		matches := regex.FindAllString(content, -1)
		urls = append(urls, matches...)
	}

	return urls, nil
}
