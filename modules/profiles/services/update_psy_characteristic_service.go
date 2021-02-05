package services

import (
	"errors"
	"strings"

	models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/identifier"
	"github.com/guicostaarantes/psi-server/utils/merge"
)

// UpdatePsyCharacteristicService is a service that edits a psychologist profile
type UpdatePsyCharacteristicService struct {
	DatabaseUtil   database.IDatabaseUtil
	IdentifierUtil identifier.IIdentifierUtil
	MergeUtil      merge.IMergeUtil
}

// Execute is the method that runs the business logic of the service
func (s UpdatePsyCharacteristicService) Execute(id string, psyCharInput *models.UpdatePsyCharacteristicInput) error {

	psyChar := models.PsyCharacteristic{}

	findErr := s.DatabaseUtil.FindOne("psi_db", "psychologist_characteristics", map[string]interface{}{"id": id}, &psyChar)
	if findErr != nil {
		return findErr
	}

	if psyChar.ID == "" {
		return errors.New("resource not found")
	}

	mergeErr := s.MergeUtil.Merge(&psyChar, psyCharInput)
	if mergeErr != nil {
		return mergeErr
	}

	psyChar.PossibleValues = strings.Join(psyCharInput.PossibleValues, ",")

	writeErr := s.DatabaseUtil.UpdateOne("psi_db", "psychologist_characteristics", map[string]interface{}{"id": id}, psyChar)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
