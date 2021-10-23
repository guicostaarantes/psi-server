package services

import (
	"fmt"

	"github.com/guicostaarantes/psi-server/modules/agreements/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// UpsertTermService is a service that creates a new term or updates an existing one
type UpsertTermService struct {
	DatabaseUtil database.IDatabaseUtil
}

func (s UpsertTermService) Execute(name string, version int64, profileType models.TermProfileType, active bool) error {

	existingTerm := models.Term{}

	findErr := s.DatabaseUtil.FindOne("terms", map[string]interface{}{"name": name, "version": version, "profileType": profileType}, &existingTerm)
	if findErr != nil {
		return findErr
	}

	if existingTerm.Name != "" {

		existingTerm.Active = active

		updateErr := s.DatabaseUtil.UpdateOne("terms", map[string]interface{}{"name": name, "version": version, "profileType": profileType}, existingTerm)
		if updateErr != nil {
			return updateErr
		}

		return nil

	}

	if version > 1 {
		previousTerm := models.Term{}

		findErr := s.DatabaseUtil.FindOne("terms", map[string]interface{}{"name": name, "version": version - 1, "profileType": profileType}, &previousTerm)
		if findErr != nil {
			return findErr
		}

		if previousTerm.Name == "" {
			return fmt.Errorf("version %d of term %s does not exist", version-1, name)
		}
	}

	newTerm := models.Term{
		Name:        name,
		Version:     version,
		ProfileType: profileType,
		Active:      active,
	}

	insertErr := s.DatabaseUtil.InsertOne("terms", newTerm)
	if insertErr != nil {
		return insertErr
	}

	return nil

}
