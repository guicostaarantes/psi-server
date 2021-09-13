package services

import (
	"io"

	files_services "github.com/guicostaarantes/psi-server/modules/files/services"
	models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/identifier"
)

// UpsertPatientService is a service that creates or updates a patient profile
type UpsertPatientService struct {
	DatabaseUtil            database.IDatabaseUtil
	IdentifierUtil          identifier.IIdentifierUtil
	UploadAvatarFileService *files_services.UploadAvatarFileService
}

// Execute is the method that runs the business logic of the service
func (s UpsertPatientService) Execute(userID string, input *models.UpsertPatientInput) error {

	existantPatient := models.Patient{}

	findErr := s.DatabaseUtil.FindOne("patients", map[string]interface{}{"userId": userID}, &existantPatient)
	if findErr != nil {
		return findErr
	}

	if existantPatient.ID != "" {
		existantPatient.FullName = input.FullName
		existantPatient.LikeName = input.LikeName
		existantPatient.BirthDate = input.BirthDate
		existantPatient.City = input.City

		if input.Avatar != nil {
			data, readErr := io.ReadAll(input.Avatar.File)
			if readErr != nil {
				return readErr
			}

			fileName, writeErr := s.UploadAvatarFileService.Execute(input.UserID, data)
			if writeErr != nil {
				return writeErr
			}

			existantPatient.Avatar = fileName
		}

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

	var avatar string

	if input.Avatar != nil {
		data, readErr := io.ReadAll(input.Avatar.File)
		if readErr != nil {
			return readErr
		}

		fileName, writeErr := s.UploadAvatarFileService.Execute(input.UserID, data)
		if writeErr != nil {
			return writeErr
		}
		avatar = fileName
	}

	patient := &models.Patient{
		ID:        patientID,
		UserID:    userID,
		FullName:  input.FullName,
		LikeName:  input.LikeName,
		BirthDate: input.BirthDate,
		City:      input.City,
		Avatar:    avatar,
	}

	writeErr := s.DatabaseUtil.InsertOne("patients", patient)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
