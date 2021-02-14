package services

import (
	"context"

	"github.com/guicostaarantes/psi-server/modules/schedule/models"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/identifier"
)

// GetPsychologistSlotsService is a service that gets all the slots of a psychologist
type GetPsychologistSlotsService struct {
	DatabaseUtil   database.IDatabaseUtil
	IdentifierUtil identifier.IIdentifierUtil
}

// Execute is the method that runs the business logic of the service
func (s GetPsychologistSlotsService) Execute(psychologistID string) ([]*models.GetPsychologistSlotsResponse, error) {

	filter := map[string]interface{}{"psychologistId": psychologistID}

	cursor, findErr := s.DatabaseUtil.FindMany("psi_db", "slots", filter)
	if findErr != nil {
		return nil, findErr
	}

	slots := []*models.GetPsychologistSlotsResponse{}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {

		slot := models.GetPsychologistSlotsResponse{}

		decodeErr := cursor.Decode(&slot)
		if decodeErr != nil {
			return nil, decodeErr
		}

		slots = append(slots, &slot)

	}

	return slots, nil

}
