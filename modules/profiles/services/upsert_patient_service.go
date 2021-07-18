package services

import (
	models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/identifier"
)

// UpsertPatientService is a service that creates or updates a patient profile
type UpsertPatientService struct {
	DatabaseUtil   database.IDatabaseUtil
	IdentifierUtil identifier.IIdentifierUtil
}

// Execute is the method that runs the business logic of the service
func (s UpsertPatientService) Execute(input *models.UpsertPatientInput) error {

	existantPatient := models.Patient{}

	findErr := s.DatabaseUtil.FindOne("patients", map[string]interface{}{"userId": input.UserID}, &existantPatient)
	if findErr != nil {
		return findErr
	}

	if existantPatient.ID != "" {
		existantPatient.FullName = input.FullName
		existantPatient.LikeName = input.LikeName
		existantPatient.BirthDate = input.BirthDate
		existantPatient.City = input.City

		writeErr := s.DatabaseUtil.UpdateOne("patients", map[string]interface{}{"id": existantPatient.ID}, existantPatient)
		if writeErr != nil {
			return writeErr
		}

		return nil
	}

	_, patientID, patientIDErr := s.IdentifierUtil.GenerateIdentifier()
	if patientIDErr != nil {
		return patientIDErr
	}

	patient := &models.Patient{
		ID:        patientID,
		UserID:    input.UserID,
		FullName:  input.FullName,
		LikeName:  input.LikeName,
		BirthDate: input.BirthDate,
		City:      input.City,
	}

	writeErr := s.DatabaseUtil.InsertOne("patients", patient)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
