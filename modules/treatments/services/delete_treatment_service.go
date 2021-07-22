package services

import (
	"errors"

	"github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// DeleteTreatmentService is a service that changes data from a treatment
type DeleteTreatmentService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s DeleteTreatmentService) Execute(id string, psychologistID string) error {

	treatment := models.Treatment{}

	findErr := s.DatabaseUtil.FindOne("treatments", map[string]interface{}{"id": id, "psychologistId": psychologistID}, &treatment)
	if findErr != nil {
		return findErr
	}

	if treatment.ID == "" {
		return errors.New("resource not found")
	}

	if treatment.Status != models.Pending {
		return errors.New("treatments can only be deleted if they their status is pending")
	}

	writeErr := s.DatabaseUtil.DeleteOne("treatments", map[string]interface{}{"id": id})
	if writeErr != nil {
		return writeErr
	}

	return nil

}
