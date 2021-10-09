package services

import (
	"io"

	files_services "github.com/guicostaarantes/psi-server/modules/files/services"
	models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/identifier"
)

// UpsertPsychologistService is a service that creates a psychologist profile
type UpsertPsychologistService struct {
	DatabaseUtil            database.IDatabaseUtil
	IdentifierUtil          identifier.IIdentifierUtil
	UploadAvatarFileService *files_services.UploadAvatarFileService
}

// Execute is the method that runs the business logic of the service
func (s UpsertPsychologistService) Execute(userID string, input *models.UpsertPsychologistInput) error {

	existantPsy := models.Psychologist{}

	findErr := s.DatabaseUtil.FindOne("psychologists", map[string]interface{}{"userId": userID}, &existantPsy)
	if findErr != nil {
		return findErr
	}

	if existantPsy.ID != "" {
		existantPsy.FullName = input.FullName
		existantPsy.LikeName = input.LikeName
		existantPsy.BirthDate = input.BirthDate
		existantPsy.City = input.City
		existantPsy.Crp = input.Crp
		existantPsy.Whatsapp = input.Whatsapp
		existantPsy.Instagram = input.Instagram
		existantPsy.Bio = input.Bio

		if input.Avatar != nil {
			data, readErr := io.ReadAll(input.Avatar.File)
			if readErr != nil {
				return readErr
			}

			fileName, writeErr := s.UploadAvatarFileService.Execute(input.UserID, data)
			if writeErr != nil {
				return writeErr
			}
			existantPsy.Avatar = fileName
		}

		writeErr := s.DatabaseUtil.UpdateOne("psychologists", map[string]interface{}{"id": existantPsy.ID}, existantPsy)
		if writeErr != nil {
			return writeErr
		}

		return nil
	}

	_, psyID, psyIDErr := s.IdentifierUtil.GenerateIdentifier()
	if psyIDErr != nil {
		return psyIDErr
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

	psy := &models.Psychologist{
		ID:        psyID,
		UserID:    userID,
		FullName:  input.FullName,
		LikeName:  input.LikeName,
		BirthDate: input.BirthDate,
		City:      input.City,
		Bio:       input.Bio,
		Avatar:    avatar,
	}

	writeErr := s.DatabaseUtil.InsertOne("psychologists", psy)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
