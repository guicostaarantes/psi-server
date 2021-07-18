package services

import (
	"context"

	"github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/identifier"
)

// GetPatientTreatmentsService is a service that gets all the treatments of a psychologist
type GetPatientTreatmentsService struct {
	DatabaseUtil   database.IDatabaseUtil
	IdentifierUtil identifier.IIdentifierUtil
}

// Execute is the method that runs the business logic of the service
func (s GetPatientTreatmentsService) Execute(patientID string) ([]*models.GetPatientTreatmentsResponse, error) {

	filter := map[string]interface{}{"patientId": patientID}

	cursor, findErr := s.DatabaseUtil.FindMany("treatments", filter)
	if findErr != nil {
		return nil, findErr
	}

	treatments := []*models.GetPatientTreatmentsResponse{}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {

		treatment := models.GetPatientTreatmentsResponse{}

		decodeErr := cursor.Decode(&treatment)
		if decodeErr != nil {
			return nil, decodeErr
		}

		treatments = append(treatments, &treatment)

	}

	return treatments, nil

}
