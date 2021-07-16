package services

import (
	"context"
	"time"

	"github.com/guicostaarantes/psi-server/modules/characteristics/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// GetTopAffinitiesForPatientService is a service that sets the top affinities for a specific patient profile if the cache is old enough, and returns them
type GetTopAffinitiesForPatientService struct {
	DatabaseUtil                      database.IDatabaseUtil
	TopAffinitiesCooldownSeconds      int64
	SetTopAffinitiesForPatientService SetTopAffinitiesForPatientService
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

		if affinity.CreatedAt+s.TopAffinitiesCooldownSeconds > time.Now().Unix() {
			topAffinities = append(topAffinities, &affinity)
		}
	}

	// if there are recent values, return them
	if len(topAffinities) > 0 {
		return topAffinities, nil
	}

	// else, run the set service to renew cache
	setErr := s.SetTopAffinitiesForPatientService.Execute(patientID)
	if setErr != nil {
		return nil, setErr
	}

	cursor, findErr = s.DatabaseUtil.FindMany("psi_db", "top_affinities", map[string]interface{}{"patientId": patientID})
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

		if affinity.CreatedAt+s.TopAffinitiesCooldownSeconds > time.Now().Unix() {
			topAffinities = append(topAffinities, &affinity)
		}
	}

	return topAffinities, nil
}
