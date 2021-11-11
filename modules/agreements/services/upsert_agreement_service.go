package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/guicostaarantes/psi-server/modules/agreements/models"
	profiles_models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/utils/identifier"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// UpsertAgreementService is a service that creates an agreement or updates an existing one
type UpsertAgreementService struct {
	IdentifierUtil identifier.IIdentifierUtil
	OrmUtil        orm.IOrmUtil
}

func (s UpsertAgreementService) Execute(profileID string, input *models.UpsertAgreementInput) error {

	var profileType models.TermProfileType

	psy := profiles_models.Psychologist{}
	pat := profiles_models.Patient{}
	result := s.OrmUtil.Db().Where("id = ?", profileID).Limit(1).Find(&pat)
	if result.Error != nil {
		return result.Error
	}
	if pat.ID != "" {
		profileType = models.Patient
	} else {
		result := s.OrmUtil.Db().Where("id = ?", profileID).Limit(1).Find(&psy)
		if result.Error != nil {
			return result.Error
		}
		if psy.ID != "" {
			profileType = models.Psychologist
		} else {
			return errors.New("resource not found")
		}
	}

	signingTerm := models.Term{}

	result = s.OrmUtil.Db().Where("name = ? AND version = ? AND profile_type = ?", input.TermName, input.TermVersion, profileType).Limit(1).Find(&signingTerm)
	if result.Error != nil {
		return result.Error
	}

	if signingTerm.Name == "" {
		return fmt.Errorf("version %d of term %s does not exist", input.TermVersion, input.TermName)
	}

	existingAgreement := models.Agreement{}

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

		existingAgreement.SignedAt = time.Now().Unix()

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

	newAgreement := models.Agreement{
		ID:          agrID,
		TermName:    input.TermName,
		TermVersion: input.TermVersion,
		ProfileID:   profileID,
		SignedAt:    time.Now().Unix(),
	}

	result = s.OrmUtil.Db().Create(&newAgreement)
	if result.Error != nil {
		return result.Error
	}

	return nil

}
