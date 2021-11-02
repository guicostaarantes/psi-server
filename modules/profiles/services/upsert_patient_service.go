package services

import (
	"io"

	files_services "github.com/guicostaarantes/psi-server/modules/files/services"
	models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/identifier"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// UpsertPatientService is a service that creates or updates a patient profile
type UpsertPatientService struct {
	DatabaseUtil            database.IDatabaseUtil
	IdentifierUtil          identifier.IIdentifierUtil
	OrmUtil                 orm.IOrmUtil
	UploadAvatarFileService *files_services.UploadAvatarFileService
}

// Execute is the method that runs the business logic of the service
func (s UpsertPatientService) Execute(userID string, input *models.UpsertPatientInput) error {

	existingPatient := models.Patient{}

	result := s.OrmUtil.Db().Where("user_id = ?", userID).Limit(1).Find(&existingPatient)
	if result.Error != nil {
		return result.Error
	}

	if existingPatient.ID != "" {
		existingPatient.FullName = input.FullName
		existingPatient.LikeName = input.LikeName
		existingPatient.BirthDate = input.BirthDate
		existingPatient.City = input.City

		if input.Avatar != nil {
			data, readErr := io.ReadAll(input.Avatar.File)
			if readErr != nil {
				return readErr
			}

			fileName, writeErr := s.UploadAvatarFileService.Execute(input.UserID, data)
			if writeErr != nil {
				return writeErr
			}

			existingPatient.Avatar = fileName
		}

		result = s.OrmUtil.Db().Save(&existingPatient)
		if result.Error != nil {
			return result.Error
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

	newPatient := models.Patient{
		ID:        patientID,
		UserID:    userID,
		FullName:  input.FullName,
		LikeName:  input.LikeName,
		BirthDate: input.BirthDate,
		City:      input.City,
		Avatar:    avatar,
	}

	result = s.OrmUtil.Db().Create(&newPatient)
	if result.Error != nil {
		return result.Error
	}

	return nil

}
