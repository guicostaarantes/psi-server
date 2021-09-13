package services

import (
	"errors"

	"github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// UpdateTreatmentService is a service that changes data from a treatment
type UpdateTreatmentService struct {
	DatabaseUtil                   database.IDatabaseUtil
	CheckTreatmentCollisionService *CheckTreatmentCollisionService
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

	if input.PriceRangeName != "" && treatment.Status == models.Pending {
		return errors.New("pending treatments are not allowed to have a price range")
	}

	if input.PriceRangeName == "" && treatment.Status != models.Pending {
		return errors.New("non-pending treatments must have a price range")
	}

	checkErr := s.CheckTreatmentCollisionService.Execute(psychologistID, input.Frequency, input.Phase, input.Duration, id)
	if checkErr != nil {
		return checkErr
	}

	treatment.Frequency = input.Frequency
	treatment.Phase = input.Phase
	treatment.Duration = input.Duration
	treatment.PriceRangeName = input.PriceRangeName

	writeErr := s.DatabaseUtil.UpdateOne("treatments", map[string]interface{}{"id": id}, treatment)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
