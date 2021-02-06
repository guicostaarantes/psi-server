package services

import (
	"errors"

	models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/identifier"
)

// CreatePsychologistService is a service that creates a psychologist profile
type CreatePsychologistService struct {
	DatabaseUtil   database.IDatabaseUtil
	IdentifierUtil identifier.IIdentifierUtil
}

// Execute is the method that runs the business logic of the service
func (s CreatePsychologistService) Execute(psyInput *models.CreatePsychologistInput) error {

	psyWithSameUserID := models.Psychologist{}

	findErr := s.DatabaseUtil.FindOne("psi_db", "psychologists", map[string]interface{}{"userId": psyInput.UserID}, &psyWithSameUserID)
	if findErr != nil {
		return findErr
	}

	if psyWithSameUserID.ID != "" {
		return errors.New("cannot create two psychologist profiles for the same user")
	}

	_, psyID, psyIDErr := s.IdentifierUtil.GenerateIdentifier()
	if psyIDErr != nil {
		return psyIDErr
	}

	psy := &models.Psychologist{
		ID:        psyID,
		UserID:    psyInput.UserID,
		BirthDate: psyInput.BirthDate,
		City:      psyInput.City,
	}

	writeErr := s.DatabaseUtil.InsertOne("psi_db", "psychologists", psy)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
