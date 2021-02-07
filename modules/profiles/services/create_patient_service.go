package services

import (
	"errors"

	models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/identifier"
)

// CreatePatientService is a service that creates a patient profile
type CreatePatientService struct {
	DatabaseUtil   database.IDatabaseUtil
	IdentifierUtil identifier.IIdentifierUtil
}

// Execute is the method that runs the business logic of the service
func (s CreatePatientService) Execute(input *models.CreatePatientInput) error {

	patientWithSameUserID := models.Patient{}

	findErr := s.DatabaseUtil.FindOne("psi_db", "patients", map[string]interface{}{"userId": input.UserID}, &patientWithSameUserID)
	if findErr != nil {
		return findErr
	}

	if patientWithSameUserID.ID != "" {
		return errors.New("cannot create two patient profiles for the same user")
	}

	_, patientID, patientIDErr := s.IdentifierUtil.GenerateIdentifier()
	if patientIDErr != nil {
		return patientIDErr
	}

	patient := &models.Patient{
		ID:        patientID,
		UserID:    input.UserID,
		BirthDate: input.BirthDate,
		City:      input.City,
	}

	writeErr := s.DatabaseUtil.InsertOne("psi_db", "patients", patient)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
