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

type EvidenceFromID struct {
	Type    EvidenceType `json:"type"`
	Remarks string       `json:"remarks"`
	Content string       `json:"content"`
	Size    string       `json:"size"`
}

type EvidenceWithID struct {
	ID      int          `json:"id"`
	Type    EvidenceType `json:"type"`
	Remarks string       `json:"remarks"`
	Content string       `json:"content"`
	Size    string       `json:"size"`
}

type ImageFromID struct {
	Type    EvidenceType `json:"type"`
	Content string       `json:"content"`
}

type HardDeleteConfirmation struct {
	Confirm string `json:"confirm"`
}
