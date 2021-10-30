package services

import (
	"context"

	"github.com/guicostaarantes/psi-server/modules/agreements/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// GetTermsByProfileTypeService is a service that gets the terms for a specific profile type
type GetTermsByProfileTypeService struct {
	DatabaseUtil database.IDatabaseUtil
}

func (s GetTermsByProfileTypeService) Execute(profileType models.TermProfileType) ([]*models.Term, error) {

	terms := []*models.Term{}

	cursor, findErr := s.DatabaseUtil.FindMany("terms", map[string]interface{}{"profileType": string(profileType)})
	if findErr != nil {
		return nil, findErr
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		term := models.Term{}

		decodeErr := cursor.Decode(&term)
		if decodeErr != nil {
			return nil, decodeErr
		}

		terms = append(terms, &term)

	}

	return terms, nil

}
