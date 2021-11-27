package agreements_services

import (
	"errors"
	"fmt"

	agreements_models "github.com/guicostaarantes/psi-server/modules/agreements/models"
	profiles_models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/utils/identifier"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// UpsertAgreementService is a service that creates an agreement or updates an existing one
type UpsertAgreementService struct {
	IdentifierUtil identifier.IIdentifierUtil
	OrmUtil        orm.IOrmUtil
}

func (s UpsertAgreementService) Execute(profileID string, input *agreements_models.UpsertAgreementInput) error {

	var profileType agreements_models.TermProfileType

	psy := profiles_models.Psychologist{}
	pat := profiles_models.Patient{}
	result := s.OrmUtil.Db().Where("id = ?", profileID).Limit(1).Find(&pat)
	if result.Error != nil {
		return result.Error
	}
	if pat.ID != "" {
		profileType = agreements_models.Patient
	} else {
		result := s.OrmUtil.Db().Where("id = ?", profileID).Limit(1).Find(&psy)
		if result.Error != nil {
			return result.Error
		}
		if psy.ID != "" {
			profileType = agreements_models.Psychologist
		} else {
			return errors.New("resource not found")
		}
	}

	signingTerm := agreements_models.Term{}

	result = s.OrmUtil.Db().Where("name = ? AND version = ? AND profile_type = ?", input.TermName, input.TermVersion, profileType).Limit(1).Find(&signingTerm)
	if result.Error != nil {
		return result.Error
	}

	if signingTerm.Name == "" {
		return fmt.Errorf("version %d of term %s does not exist", input.TermVersion, input.TermName)
	}

	existingAgreement := agreements_models.Agreement{}

	result = s.OrmUtil.Db().Where("term_name = ? AND term_version = ? AND profile_id = ?", input.TermName, input.TermVersion, profileID).Limit(1).Find(&existingAgreement)
	if result.Error != nil {
		return result.Error
	}

	if !input.Agreed {

		if existingAgreement.ID != "" {
			result = s.OrmUtil.Db().Delete(&existingAgreement)
			if result.Error != nil {
				return result.Error
			}
		}

		return nil

	}

	if existingAgreement.ID != "" {

		result = s.OrmUtil.Db().Save(&existingAgreement)
		if result.Error != nil {
			return result.Error
		}

		return nil

	}

	_, agrID, agrIDErr := s.IdentifierUtil.GenerateIdentifier()
	if agrIDErr != nil {
		return agrIDErr
	}

	newAgreement := agreements_models.Agreement{
		ID:          agrID,
		TermName:    input.TermName,
		TermVersion: input.TermVersion,
		ProfileID:   profileID,
	}

	result = s.OrmUtil.Db().Create(&newAgreement)
	if result.Error != nil {
		return result.Error
	}

	return nil

}
