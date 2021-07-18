package services

import (
	models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/identifier"
)

// UpsertPsychologistService is a service that creates a psychologist profile
type UpsertPsychologistService struct {
	DatabaseUtil   database.IDatabaseUtil
	IdentifierUtil identifier.IIdentifierUtil
}

// Execute is the method that runs the business logic of the service
func (s UpsertPsychologistService) Execute(input *models.UpsertPsychologistInput) error {

	existantPsy := models.Psychologist{}

	findErr := s.DatabaseUtil.FindOne("psychologists", map[string]interface{}{"userId": input.UserID}, &existantPsy)
	if findErr != nil {
		return findErr
	}

	if existantPsy.ID != "" {
		existantPsy.FullName = input.FullName
		existantPsy.LikeName = input.LikeName
		existantPsy.BirthDate = input.BirthDate
		existantPsy.City = input.City

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

	psy := &models.Psychologist{
		ID:        psyID,
		UserID:    input.UserID,
		FullName:  input.FullName,
		LikeName:  input.LikeName,
		BirthDate: input.BirthDate,
		City:      input.City,
	}

	writeErr := s.DatabaseUtil.InsertOne("psychologists", psy)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
