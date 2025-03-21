package models

type User struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Password       string `json:"-"`
	Role           string `json:"role"`
	ClearanceLevel string `json:"clearance_level"`
	Deleted        bool   `json:"deleted"`
}
