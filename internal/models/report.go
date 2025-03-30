package models

type CrimeReportRequest struct {
	Email       string `json:"email"`
	CivilID     string `json:"civil_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Area        string `json:"area"`
	City        string `json:"city"`
}

type Report struct {
	ReportID    int     `json:"report_id"`
	Email       string  `json:"email"`
	CivilID     string  `json:"civil_id"`
	Name        string  `json:"name"`
	Role        string  `json:"role"`
	CaseNumber  *string `json:"case_number,omitempty"`
	Description string  `json:"description"`
	Area        string  `json:"area"`
	City        string  `json:"city"`
}
