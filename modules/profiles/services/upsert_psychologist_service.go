package services

import (
	"io"

	files_services "github.com/guicostaarantes/psi-server/modules/files/services"
	profiles_models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/utils/identifier"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// UpsertPsychologistService is a service that creates a psychologist profile
type UpsertPsychologistService struct {
	IdentifierUtil          identifier.IIdentifierUtil
	OrmUtil                 orm.IOrmUtil
	UploadAvatarFileService *files_services.UploadAvatarFileService
}

// Execute is the method that runs the business logic of the service
func (s UpsertPsychologistService) Execute(userID string, input *profiles_models.UpsertPsychologistInput) error {

	existingPsy := profiles_models.Psychologist{}

	result := s.OrmUtil.Db().Where("user_id = ?", userID).Limit(1).Find(&existingPsy)
	if result.Error != nil {
		return result.Error
	}

	if existingPsy.ID != "" {
		existingPsy.FullName = input.FullName
		existingPsy.LikeName = input.LikeName
		existingPsy.BirthDate = input.BirthDate
		existingPsy.City = input.City
		existingPsy.Crp = input.Crp
		existingPsy.Whatsapp = input.Whatsapp
		existingPsy.Instagram = input.Instagram
		existingPsy.Bio = input.Bio

		if input.Avatar != nil {
			data, readErr := io.ReadAll(input.Avatar.File)
			if readErr != nil {
				return readErr
			}

			fileName, writeErr := s.UploadAvatarFileService.Execute(input.UserID, data)
			if writeErr != nil {
				return writeErr
			}
			existingPsy.Avatar = fileName
		}

		result = s.OrmUtil.Db().Save(&existingPsy)
		if result.Error != nil {
			return result.Error
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

	newPsy := profiles_models.Psychologist{
		ID:        psyID,
		UserID:    userID,
		FullName:  input.FullName,
		LikeName:  input.LikeName,
		BirthDate: input.BirthDate,
		City:      input.City,
		Crp:       input.Crp,
		Whatsapp:  input.Whatsapp,
		Instagram: input.Instagram,
		Bio:       input.Bio,
		Avatar:    avatar,
	}

	result = s.OrmUtil.Db().Create(&newPsy)
	if result.Error != nil {
		return result.Error
	}

	return nil

}
