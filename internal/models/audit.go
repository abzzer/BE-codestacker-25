package models

type AuditAction string

const (
	AuditAdded       AuditAction = "added"
	AuditUpdated     AuditAction = "updated"
	AuditSoftDeleted AuditAction = "soft_deleted"
	AuditHardDeleted AuditAction = "hard_deleted"
)

type AuditLog struct {
	ID         int         `json:"id"`
	Action     AuditAction `json:"action"`
	EvidenceID int         `json:"evidence_id"`
	UserID     string      `json:"user_id"`
	Timestamp  string      `json:"timestamp"`
}
