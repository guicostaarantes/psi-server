package services

import (
	"context"
	"time"

	"github.com/guicostaarantes/psi-server/modules/characteristics/models"
	cooldowns_models "github.com/guicostaarantes/psi-server/modules/cooldowns/models"
	cooldowns_services "github.com/guicostaarantes/psi-server/modules/cooldowns/services"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// GetTopAffinitiesForPatientService is a service that sets the top affinities for a specific patient profile if the cache is old enough, and returns them
type GetTopAffinitiesForPatientService struct {
	DatabaseUtil                      database.IDatabaseUtil
	TopAffinitiesCooldownSeconds      int64
	GetCooldownService                *cooldowns_services.GetCooldownService
	SetTopAffinitiesForPatientService *SetTopAffinitiesForPatientService
}

// Execute is the method that runs the business logic of the service
func (s GetTopAffinitiesForPatientService) Execute(patientID string) ([]*models.Affinity, error) {
	cooldown, getErr := s.GetCooldownService.Execute(patientID, cooldowns_models.Patient, cooldowns_models.TopAffinitiesSet)
	if getErr != nil {
		return nil, getErr
	}

	if cooldown == nil {
		setErr := s.SetTopAffinitiesForPatientService.Execute(patientID)
		if setErr != nil {
			return nil, setErr
		}
	}

	topAffinities := []*models.Affinity{}

	cursor, findErr := s.DatabaseUtil.FindMany("top_affinities", map[string]interface{}{"patientId": patientID})
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
