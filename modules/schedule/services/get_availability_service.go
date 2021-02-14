package services

import (
	"context"

	"github.com/guicostaarantes/psi-server/modules/schedule/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// GetAvailabilityService is a service that gets all the availabities for the user within a specific period
type GetAvailabilityService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s GetAvailabilityService) Execute(id string) ([]*models.AvailabilityResponse, error) {

	availabilities := []*models.AvailabilityResponse{}

	cursor, findErr := s.DatabaseUtil.FindMany("psi_db", "availabilities", map[string]interface{}{"psychologistId": id})
	if findErr != nil {
		return nil, findErr
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		availability := models.AvailabilityResponse{}

		decodeErr := cursor.Decode(&availability)
		if decodeErr != nil {
			return nil, decodeErr
		}

		availabilities = append(availabilities, &availability)

	}

	return availabilities, nil

}
