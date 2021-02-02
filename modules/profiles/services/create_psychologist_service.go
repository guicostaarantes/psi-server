package profiles_services

import (
	"errors"

	models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/identifier"
	"github.com/guicostaarantes/psi-server/utils/merge"
)

// CreatePsychologistService is a service that creates a psychologist profile
type CreatePsychologistService struct {
	DatabaseUtil   database.IDatabaseUtil
	IdentifierUtil identifier.IIdentifierUtil
	MergeUtil      merge.IMergeUtil
}

// Execute is the method that runs the business logic of the service
func (s CreatePsychologistService) Execute(psyInput *models.CreatePsychologistInput) error {

	psyWithSameUserID := models.Psychologist{}

	findErr := s.DatabaseUtil.FindOne("psi_db", "psychologists", "userId", psyInput.UserID, &psyWithSameUserID)
	if findErr != nil {
		return findErr
	}

	if psyWithSameUserID.ID != "" {
		return errors.New("cannot create two psychologist profiles for the same user")
	}

	_, psyID, psyIdErr := s.IdentifierUtil.GenerateIdentifier()
	if psyIdErr != nil {
		return psyIdErr
	}

	psy := &models.Psychologist{
		ID: psyID,
	}

	mergeErr := s.MergeUtil.Merge(&psy, psyInput)
	if mergeErr != nil {
		return mergeErr
	}

	writeErr := s.DatabaseUtil.InsertOne("psi_db", "psychologists", psy)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
