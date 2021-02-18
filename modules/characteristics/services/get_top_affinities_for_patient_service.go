package services

import (
	"context"

	"github.com/guicostaarantes/psi-server/modules/characteristics/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// GetTopAffinitiesForPatientService is a service that gets all top affinities for a specific patient profile
type GetTopAffinitiesForPatientService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s GetTopAffinitiesForPatientService) Execute(patientID string) ([]*models.Affinity, error) {

	topAffinities := []*models.Affinity{}

	cursor, findErr := s.DatabaseUtil.FindMany("psi_db", "top_affinities", map[string]interface{}{"patientId": patientID})
	if findErr != nil {
		return nil, findErr
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		affinity := models.Affinity{}

		decodeErr := cursor.Decode(&affinity)
		if decodeErr != nil {
			return nil, decodeErr
		}

		topAffinities = append(topAffinities, &affinity)
	}

	return topAffinities, nil

}
