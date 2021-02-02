package profiles_services

import (
	"errors"
	"strings"

	models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/identifier"
	"github.com/guicostaarantes/psi-server/utils/merge"
)

// CreatePsyCharacteristicService is a service that creates a psychologist profile
type CreatePsyCharacteristicService struct {
	DatabaseUtil   database.IDatabaseUtil
	IdentifierUtil identifier.IIdentifierUtil
	MergeUtil      merge.IMergeUtil
}

// Execute is the method that runs the business logic of the service
func (s CreatePsyCharacteristicService) Execute(psyCharInput *models.CreatePsyCharacteristicInput) error {

	psyCharWithSameName := models.PsyCharacteristic{}

	findErr := s.DatabaseUtil.FindOne("psi_db", "psychologist_characteristics", "name", psyCharInput.Name, &psyCharWithSameName)
	if findErr != nil {
		return findErr
	}

	if psyCharWithSameName.ID != "" {
		return errors.New("cannot create two psychologist characteristics with the same name")
	}

	_, psyID, psyIdErr := s.IdentifierUtil.GenerateIdentifier()
	if psyIdErr != nil {
		return psyIdErr
	}

	psyChar := &models.PsyCharacteristic{
		ID: psyID,
	}

	mergeErr := s.MergeUtil.Merge(&psyChar, psyCharInput)
	if mergeErr != nil {
		return mergeErr
	}

	psyChar.Values = strings.Join(psyCharInput.Values, ",")

	writeErr := s.DatabaseUtil.InsertOne("psi_db", "psychologist_characteristics", psyChar)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
