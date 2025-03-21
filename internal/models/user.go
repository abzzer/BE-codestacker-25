package models

type User struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Password       string `json:"-"`
	Role           string `json:"role"`
	ClearanceLevel string `json:"clearance_level"`
	Deleted        bool   `json:"-"`
}

type UpdateUserInput struct {
	Name           *string `json:"name"`
	Password       *string `json:"password"`
	Role           *string `json:"role"`
	ClearanceLevel *string `json:"clearance_level"`
}

type LoginRequest struct {
	UserID   string `json:"user_id"`
	Password string `json:"password"`
}
