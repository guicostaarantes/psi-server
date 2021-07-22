package services

import (
	"context"
	"errors"

	"github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// CheckTreatmentCollisionService is a service that checks if a treatment period collides with others from the same psychologist
type CheckTreatmentCollisionService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s CheckTreatmentCollisionService) Execute(psychologistID string, weeklyStart int64, duration int64, updatingID string) error {

	weekDuration := int64(7 * 24 * 60 * 60)

	end := (weeklyStart + duration) % weekDuration

	cursor, findErr := s.DatabaseUtil.FindMany("treatments", map[string]interface{}{"psychologistId": psychologistID})
	if findErr != nil {
		return findErr
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		treatment := models.Treatment{}

		decodeErr := cursor.Decode(&treatment)
		if decodeErr != nil {
			return decodeErr
		}

		treatmentEnd := (treatment.WeeklyStart + treatment.Duration) % weekDuration

		// If 3 of the 4 conditions below are true, it means there is no clash between treatments
		noClash := 0
		if weeklyStart < end {
			noClash++
		}
		if end <= treatment.WeeklyStart {
			noClash++
		}
		if treatment.WeeklyStart < treatmentEnd {
			noClash++
		}
		if treatmentEnd <= weeklyStart {
			noClash++
		}

		if noClash < 3 && treatment.ID != updatingID {
			return errors.New("there is another treatment in the same period")
		}
	}

	return nil

}
