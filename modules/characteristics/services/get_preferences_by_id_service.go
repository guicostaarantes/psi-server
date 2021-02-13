package services

import (
	"context"

	"github.com/guicostaarantes/psi-server/modules/characteristics/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// GetPreferencesByIDService is a service that gets the preferences of a profile based on its id
type GetPreferencesByIDService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s GetPreferencesByIDService) Execute(id string) ([]*models.PreferenceResponse, error) {

	preferences := []*models.PreferenceResponse{}

	cursor, findErr := s.DatabaseUtil.FindMany("psi_db", "preferences", map[string]interface{}{"profileId": id})
	if findErr != nil {
		return nil, findErr
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		preference := models.PreferenceResponse{}

		decodeErr := cursor.Decode(&preference)
		if decodeErr != nil {
			return nil, decodeErr
		}

		preferences = append(preferences, &preference)
	}

	return preferences, nil

}
