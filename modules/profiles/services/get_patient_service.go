package services

import (
	models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// GetPatientService is a service that gets the patient profile based on id
type GetPatientService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s GetPatientService) Execute(id string) (*models.Patient, error) {

	patient := &models.Patient{}

	findErr := s.DatabaseUtil.FindOne("patients", map[string]interface{}{"id": id}, patient)
	if findErr != nil {
		return nil, findErr
	}

	if patient.ID == "" {
		return nil, nil
	}

	return patient, nil

}
