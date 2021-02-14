package services

import (
	"context"

	"github.com/guicostaarantes/psi-server/modules/schedule/models"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/identifier"
)

// GetPatientSlotsService is a service that gets all the slots of a psychologist
type GetPatientSlotsService struct {
	DatabaseUtil   database.IDatabaseUtil
	IdentifierUtil identifier.IIdentifierUtil
}

// Execute is the method that runs the business logic of the service
func (s GetPatientSlotsService) Execute(patientID string) ([]*models.GetPatientSlotsResponse, error) {

	filter := map[string]interface{}{"patientId": patientID}

	cursor, findErr := s.DatabaseUtil.FindMany("psi_db", "slots", filter)
	if findErr != nil {
		return nil, findErr
	}

	slots := []*models.GetPatientSlotsResponse{}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {

		slot := models.GetPatientSlotsResponse{}

		decodeErr := cursor.Decode(&slot)
		if decodeErr != nil {
			return nil, decodeErr
		}

		slots = append(slots, &slot)

	}

	return slots, nil

}
