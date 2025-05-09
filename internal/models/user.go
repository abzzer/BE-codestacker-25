package models

type UserRole string
type ClearanceLevel string

const (
	RoleAdmin        UserRole = "admin"
	RoleInvestigator UserRole = "investigator"
	RoleOfficer      UserRole = "officer"
	RoleAuditor      UserRole = "auditor"

	ClearanceLow      ClearanceLevel = "low"
	ClearanceMedium   ClearanceLevel = "medium"
	ClearanceHigh     ClearanceLevel = "high"
	ClearanceCritical ClearanceLevel = "critical"
)

type LoginRequest struct {
	UserID   string `json:"user_id"`
	Password string `json:"password"`
}

type User struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Password       string         `json:"-"`
	Role           UserRole       `json:"role"`
	ClearanceLevel ClearanceLevel `json:"clearance_level"`
	Deleted        bool           `json:"-"`
}

type UserCreate struct {
	Name           string         `json:"name"`
	Password       string         `json:"password"`
	Role           UserRole       `json:"role"`
	ClearanceLevel ClearanceLevel `json:"clearance_level"`
}

type UpdateUser struct {
	ID             string          `json:"-"`
	Name           *string         `json:"name"`
	Password       *string         `json:"password"`
	Role           *UserRole       `json:"role"`
	ClearanceLevel *ClearanceLevel `json:"clearance_level"`
}
