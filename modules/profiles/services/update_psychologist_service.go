package profiles_services

import (
	"errors"

	models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/merge"
)

// UpdatePsychologistService is a service that edits a psychologist profile
type UpdatePsychologistService struct {
	DatabaseUtil database.IDatabaseUtil
	MergeUtil    merge.IMergeUtil
}

// Execute is the method that runs the business logic of the service
func (s UpdatePsychologistService) Execute(id string, psyInput *models.UpdatePsychologistInput) error {

	psy := models.Psychologist{}

	findErr := s.DatabaseUtil.FindOne("psi_db", "psychologists", "id", id, &psy)
	if findErr != nil {
		return findErr
	}

	if psy.ID == "" {
		return errors.New("resource not found")
	}

	mergeErr := s.MergeUtil.Merge(&psy, psyInput)
	if mergeErr != nil {
		return mergeErr
	}

	writeErr := s.DatabaseUtil.UpdateOne("psi_db", "psychologists", "id", id, psy)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
