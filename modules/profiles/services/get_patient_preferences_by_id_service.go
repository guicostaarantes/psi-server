package services

import (
	"context"

	"github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// GetPatientPreferencesByPatientIDService is a service that gets the preferences of a patient by patient id
type GetPatientPreferencesByPatientIDService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s GetPatientPreferencesByPatientIDService) Execute(id string) ([]*models.PatientPreference, error) {

	preferences := []*models.PatientPreference{}

	cursor, findErr := s.DatabaseUtil.FindMany("psi_db", "patient_preferences", map[string]interface{}{"patientId": id})
	if findErr != nil {
		return nil, findErr
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		preference := models.PatientPreference{}

		decodeErr := cursor.Decode(&preference)
		if decodeErr != nil {
			return nil, decodeErr
		}

		preferences = append(preferences, &preference)
	}

	return preferences, nil

}
