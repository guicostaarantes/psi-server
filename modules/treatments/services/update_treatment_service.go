package services

import (
	"errors"

	"github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// UpdateTreatmentService is a service that changes data from a treatment
type UpdateTreatmentService struct {
	DatabaseUtil            database.IDatabaseUtil
	ScheduleIntervalSeconds int64
}

// Execute is the method that runs the business logic of the service
func (s UpdateTreatmentService) Execute(id string, psychologistID string, input models.UpdateTreatmentInput) error {

	treatment := models.Treatment{}

	findErr := s.DatabaseUtil.FindOne("treatments", map[string]interface{}{"id": id, "psychologistId": psychologistID}, &treatment)
	if findErr != nil {
		return findErr
	}

	if treatment.ID == "" {
		return errors.New("resource not found")
	}

	if input.PriceRange != "" && treatment.Status == models.Pending {
		return errors.New("pending treatments are not allowed to have a price range")
	}

	if input.PriceRange == "" && treatment.Status != models.Pending {
		return errors.New("non-pending treatments must have a price range")
	}

	checkTreatmentCollisionService := CheckTreatmentCollisionService{
		DatabaseUtil:            s.DatabaseUtil,
		ScheduleIntervalSeconds: s.ScheduleIntervalSeconds,
	}

	checkErr := checkTreatmentCollisionService.Execute(psychologistID, input.Frequency, input.Phase, input.Duration, id)
	if checkErr != nil {
		return checkErr
	}

	treatment.Frequency = input.Frequency
	treatment.Phase = input.Phase
	treatment.Duration = input.Duration
	treatment.PriceRange = input.PriceRange

	writeErr := s.DatabaseUtil.UpdateOne("treatments", map[string]interface{}{"id": id}, treatment)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
