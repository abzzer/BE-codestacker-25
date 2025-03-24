package repository

import (
	"errors"

	"github.com/abzzer/BE-codestacker-25/internal/models"
)

var clearanceOrder = map[models.CaseLevel]int{models.Low: 1, models.Medium: 2, models.High: 3, models.Critical: 4}

func IsClearanceSufficient(userClearance, required models.CaseLevel) bool {
	return clearanceOrder[userClearance] >= clearanceOrder[required]
}

func CheckClearance(userRole, userClearance string, caseLevel models.CaseLevel) error {
	if userRole == "admin" || userRole == "investigator" {
		return nil
	}

	userLevel := models.CaseLevel(userClearance)
	if !IsClearanceSufficient(userLevel, caseLevel) {
		return errors.New("insufficient clearance level for this case")
	}

	return nil
}
