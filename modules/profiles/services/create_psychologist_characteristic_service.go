package services

import (
	"errors"
	"strings"

	models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/identifier"
)

// CreatePsychologistCharacteristicService is a service that creates a psychologist profile
type CreatePsychologistCharacteristicService struct {
	DatabaseUtil   database.IDatabaseUtil
	IdentifierUtil identifier.IIdentifierUtil
}

// Execute is the method that runs the business logic of the service
func (s CreatePsychologistCharacteristicService) Execute(psyCharInput *models.CreatePsychologistCharacteristicInput) error {

	psyCharWithSameName := models.PsychologistCharacteristic{}

	findErr := s.DatabaseUtil.FindOne("psi_db", "psychologist_characteristics", map[string]interface{}{"name": psyCharInput.Name}, &psyCharWithSameName)
	if findErr != nil {
		return findErr
	}

	if psyCharWithSameName.ID != "" {
		return errors.New("cannot create two psychologist characteristics with the same name")
	}

	_, psyID, psyIDErr := s.IdentifierUtil.GenerateIdentifier()
	if psyIDErr != nil {
		return psyIDErr
	}

	psyChar := &models.PsychologistCharacteristic{
		ID:             psyID,
		Name:           psyCharInput.Name,
		Many:           psyCharInput.Many,
		PossibleValues: strings.Join(psyCharInput.PossibleValues, ","),
	}

	writeErr := s.DatabaseUtil.InsertOne("psi_db", "psychologist_characteristics", psyChar)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
