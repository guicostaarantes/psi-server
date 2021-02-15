package services

import (
	"errors"

	"github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// UpdateTreatmentService is a service that changes data from a treatment
type UpdateTreatmentService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s UpdateTreatmentService) Execute(id string, psychologistID string, input models.UpdateTreatmentInput) error {

	treatment := models.Treatment{}

	findErr := s.DatabaseUtil.FindOne("psi_db", "treatments", map[string]interface{}{"id": id, "psychologistId": psychologistID}, &treatment)
	if findErr != nil {
		return findErr
	}

	if treatment.ID == "" {
		return errors.New("resource not found")
	}

	treatment.Duration = input.Duration
	treatment.Price = input.Price
	treatment.Interval = input.Interval

	writeErr := s.DatabaseUtil.UpdateOne("psi_db", "treatments", map[string]interface{}{"id": id}, treatment)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
