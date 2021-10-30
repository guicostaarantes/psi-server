package services

import (
	"context"

	"github.com/guicostaarantes/psi-server/modules/agreements/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// GetAgreementsByProfileIdService is a service that gets the agreements for a specific profileId
type GetAgreementsByProfileIdService struct {
	DatabaseUtil database.IDatabaseUtil
}

func (s GetAgreementsByProfileIdService) Execute(profileId string, profileType models.TermProfileType) ([]*models.Agreement, error) {

	agreements := []*models.Agreement{}

	cursor, findErr := s.DatabaseUtil.FindMany("agreements", map[string]interface{}{"profileId": profileId})
	if findErr != nil {
		return nil, findErr
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		agreement := models.Agreement{}

		decodeErr := cursor.Decode(&agreement)
		if decodeErr != nil {
			return nil, decodeErr
		}

		agreements = append(agreements, &agreement)

	}

	return agreements, nil

}
