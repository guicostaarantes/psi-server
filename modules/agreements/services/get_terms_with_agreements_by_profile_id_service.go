package services

import (
	"context"

	"github.com/guicostaarantes/psi-server/modules/agreements/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// GetTermsWithAgreementsByProfileIdService is a service that gets the most current versions of all active terms and their respective agreements for a specific profileId
type GetTermsWithAgreementsByProfileIdService struct {
	DatabaseUtil database.IDatabaseUtil
}

func (s GetTermsWithAgreementsByProfileIdService) Execute(profileId string, profileType models.TermProfileType) ([]*models.TermWithAgreement, error) {

	termsWithAgreements := []*models.TermWithAgreement{}

	termCursor, findErr := s.DatabaseUtil.FindMany("terms", map[string]interface{}{"profileType": string(profileType), "active": true})
	if findErr != nil {
		return nil, findErr
	}

	defer termCursor.Close(context.Background())

	for termCursor.Next(context.Background()) {
		term := models.Term{}

		decodeErr := termCursor.Decode(&term)
		if decodeErr != nil {
			return nil, decodeErr
		}

		replaced := false

		for _, v := range termsWithAgreements {
			if term.Name == v.Term.Name && term.Version > v.Term.Version {
				replaced = true
				v.Term = &term
				break
			}
		}

		if !replaced {
			termsWithAgreements = append(termsWithAgreements, &models.TermWithAgreement{Term: &term})
		}
	}

	agreementCursor, findErr := s.DatabaseUtil.FindMany("agreements", map[string]interface{}{"profileId": profileId})
	if findErr != nil {
		return nil, findErr
	}

	defer agreementCursor.Close(context.Background())

	for agreementCursor.Next(context.Background()) {
		agreement := models.Agreement{}

		decodeErr := agreementCursor.Decode(&agreement)
		if decodeErr != nil {
			return nil, decodeErr
		}

		for _, v := range termsWithAgreements {
			if v.Term.Name == agreement.TermName && v.Term.Version == agreement.TermVersion {
				v.Agreement = &agreement
				break
			}
		}

	}

	return termsWithAgreements, nil

}
