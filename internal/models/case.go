package models

import "time"

// DON'T FORGET OFFICER USE CASES!! COME BACK TO REVIEW THEM LATER!! -> Like they can change leve beacsuse officer
// also update the reports related to the specific case -> Ensure we move them up and then also link them to case

type CaseLevel string
type CaseStatus string
type PersonType string
type Gender string

const (
	Low      CaseLevel = "low"
	Medium   CaseLevel = "medium"
	High     CaseLevel = "high"
	Critical CaseLevel = "critical"

	StatusPending CaseStatus = "pending"
	StatusOngoing CaseStatus = "ongoing"
	StatusClosed  CaseStatus = "closed"

	PersonVictim  PersonType = "victim"
	PersonSuspect PersonType = "suspect"
	PersonWitness PersonType = "witness"

	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

type Case struct {
	CaseNumber  string     `json:"case_number"`
	CaseName    string     `json:"case_name"`
	Description string     `json:"description"`
	Area        string     `json:"area"`
	City        string     `json:"city"`
	CreatedBy   string     `json:"created_by"`
	CaseType    string     `json:"case_type"`
	Level       CaseLevel  `json:"level"`
	Status      CaseStatus `json:"status"`
	CreatedAt   string     `json:"created_at"`
}

type CaseRequest struct {
	CaseName    string    `json:"case_name"`
	Description string    `json:"description"`
	Area        string    `json:"area"`
	City        string    `json:"city"`
	Level       CaseLevel `json:"level"`
}

type CaseUpdate struct {
	CaseNumber  string      `json:"-"`
	CaseName    *string     `json:"case_name"`
	Description *string     `json:"description"`
	Area        *string     `json:"area"`
	City        *string     `json:"city"`
	Level       *CaseLevel  `json:"level"`
	Status      *CaseStatus `json:"status"`
}

type CaseStatusUpdate struct {
	CaseNumber string     `json:"-"`
	Status     CaseStatus `json:"status"`
}

type Person struct {
	ID         int        `json:"id"`
	CaseNumber string     `json:"case_number"`
	Type       PersonType `json:"type"`
	Name       string     `json:"name"`
	Age        int        `json:"age"`
	Gender     Gender     `json:"gender"`
	Role       string     `json:"role"`
}

type PersonRequest struct {
	CaseNumber string     `json:"-"`
	Type       PersonType `json:"type"`
	Name       string     `json:"name"`
	Age        int        `json:"age"`
	Gender     Gender     `json:"gender"`
	Role       string     `json:"role"`
}

type CaseDetailsResponse struct {
	CaseNumber   string     `json:"case_number"`
	CaseName     string     `json:"case_name"`
	Description  string     `json:"description"`
	Area         string     `json:"area"`
	City         string     `json:"city"`
	CreatedBy    string     `json:"created_by"`
	CreatedAt    time.Time  `json:"created_at"`
	CaseType     string     `json:"case_type"`
	Level        CaseLevel  `json:"level"`
	Status       CaseStatus `json:"status"`
	ReportedBy   int        `json:"reported_by"`
	NumAssignees int        `json:"num_assignees"`
	NumEvidences int        `json:"num_evidences"`
	NumSuspects  int        `json:"num_suspects"`
	NumVictims   int        `json:"num_victims"`
	NumWitnesses int        `json:"num_witnesses"`
}

type FullCaseDetails struct {
	CaseDetailsResponse
	Assignees []User           `json:"assignees"`
	Evidence  []EvidenceWithID `json:"evidence"`
	People    []Person         `json:"people"`
}
