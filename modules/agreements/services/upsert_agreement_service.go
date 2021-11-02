package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/guicostaarantes/psi-server/modules/agreements/models"
	profiles_models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/identifier"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// UpsertAgreementService is a service that creates an agreement or updates an existing one
type UpsertAgreementService struct {
	DatabaseUtil   database.IDatabaseUtil
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

	findErr := s.DatabaseUtil.FindOne("terms", map[string]interface{}{"name": input.TermName, "version": float64(input.TermVersion), "profileType": string(profileType)}, &signingTerm)
	if findErr != nil {
		return findErr
	}

	if signingTerm.Name == "" {
		return fmt.Errorf("version %d of term %s does not exist", input.TermVersion, input.TermName)
	}

	existingAgreement := models.Agreement{}

	findErr = s.DatabaseUtil.FindOne("agreements", map[string]interface{}{"termName": input.TermName, "termVersion": float64(input.TermVersion), "profileId": profileID}, &existingAgreement)
	if findErr != nil {
		return findErr
	}

	if !input.Agreed {

		if existingAgreement.ID != "" {
			deleteErr := s.DatabaseUtil.DeleteOne("agreements", map[string]interface{}{"termName": input.TermName, "termVersion": float64(input.TermVersion), "profileId": profileID})
			if deleteErr != nil {
				return deleteErr
			}
		}

		return nil

	}

	if existingAgreement.ID != "" {

		existingAgreement.SignedAt = time.Now().Unix()

		updateErr := s.DatabaseUtil.UpdateOne("agreements", map[string]interface{}{"termName": input.TermName, "termVersion": float64(input.TermVersion), "profileId": profileID}, existingAgreement)
		if updateErr != nil {
			return updateErr
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

	insertErr := s.DatabaseUtil.InsertOne("agreements", newAgreement)
	if insertErr != nil {
		return insertErr
	}

	return nil

}
