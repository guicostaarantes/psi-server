package services

import (
	"fmt"

	"github.com/guicostaarantes/psi-server/modules/agreements/models"
	"github.com/guicostaarantes/psi-server/utils/identifier"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// UpsertTermService is a service that creates a new term or updates an existing one
type UpsertTermService struct {
	IdentifierUtil identifier.IIdentifierUtil
	OrmUtil        orm.IOrmUtil
}

func (s UpsertTermService) Execute(name string, version int64, profileType models.TermProfileType, active bool) error {

	existingTerm := models.Term{}

	result := s.OrmUtil.Db().Where("name = ? AND version = ? AND profile_type = ?", name, version, profileType).Limit(1).Find(&existingTerm)
	if result.Error != nil {
		return result.Error
	}

	if existingTerm.Name != "" {

		existingTerm.Active = active

		result = s.OrmUtil.Db().Save(&existingTerm)
		if result.Error != nil {
			return result.Error
		}

		return nil

	}

	if version > 1 {
		previousTerm := models.Term{}

		result := s.OrmUtil.Db().Where("name = ? AND version = ? AND profile_type = ?", name, version-1, profileType).Limit(1).Find(&previousTerm)
		if result.Error != nil {
			return result.Error
		}

		if previousTerm.Name == "" {
			return fmt.Errorf("version %d of term %s does not exist", version-1, name)
		}
	}

	_, termID, termIDErr := s.IdentifierUtil.GenerateIdentifier()
	if termIDErr != nil {
		return termIDErr
	}

	newTerm := models.Term{
		ID:          termID,
		Name:        name,
		Version:     version,
		ProfileType: profileType,
		Active:      active,
	}

	result = s.OrmUtil.Db().Create(&newTerm)
	if result.Error != nil {
		return result.Error
	}

	return nil

}
