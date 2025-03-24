package models

type EvidenceType string

const (
	EvidenceText  EvidenceType = "text"
	EvidenceImage EvidenceType = "image"
)

type EvidenceTextRequest struct {
	CaseNumber string       `json:"case_number"`
	OfficerID  string       `json:"-"`
	Type       EvidenceType `json:"type"`
	Content    string       `json:"content"` // MinIO URL or text
	Remarks    string       `json:"remarks"`
	Size       string       `json:"size"`
}
