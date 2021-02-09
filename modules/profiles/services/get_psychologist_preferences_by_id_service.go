package services

import (
	"context"

	"github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// GetPsychologistPreferencesByPsyIDService is a service that gets the preferences of a psychologist by psychologist id
type GetPsychologistPreferencesByPsyIDService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s GetPsychologistPreferencesByPsyIDService) Execute(id string) ([]*models.PsychologistPreference, error) {

	preferences := []*models.PsychologistPreference{}

	cursor, findErr := s.DatabaseUtil.FindMany("psi_db", "psychologist_preferences", map[string]interface{}{"psychologistId": id})
	if findErr != nil {
		return nil, findErr
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		preference := models.PsychologistPreference{}

		decodeErr := cursor.Decode(&preference)
		if decodeErr != nil {
			return nil, decodeErr
		}

		preferences = append(preferences, &preference)
	}

	return preferences, nil

}
