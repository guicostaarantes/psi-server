package services

import (
	"errors"

	models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// UpdatePatientService is a service that edits a patient profile
type UpdatePatientService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s UpdatePatientService) Execute(id string, psyInput *models.UpdatePatientInput) error {

	patient := models.Patient{}

	findErr := s.DatabaseUtil.FindOne("psi_db", "patients", map[string]interface{}{"id": id}, &patient)
	if findErr != nil {
		return findErr
	}

	if patient.ID == "" {
		return errors.New("resource not found")
	}

	patient.BirthDate = psyInput.BirthDate
	patient.City = psyInput.City

	writeErr := s.DatabaseUtil.UpdateOne("psi_db", "patients", map[string]interface{}{"id": id}, patient)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
