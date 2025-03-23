package models

type CrimeReportRequest struct {
	Email       string `json:"email"`
	CivilID     string `json:"civil_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Area        string `json:"area"`
	City        string `json:"city"`
}
