package services

import (
	"context"

	"github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/identifier"
)

// GetPsychologistTreatmentsService is a service that gets all the treatments of a psychologist
type GetPsychologistTreatmentsService struct {
	DatabaseUtil   database.IDatabaseUtil
	IdentifierUtil identifier.IIdentifierUtil
}

// Execute is the method that runs the business logic of the service
func (s GetPsychologistTreatmentsService) Execute(psychologistID string) ([]*models.GetPsychologistTreatmentsResponse, error) {

	filter := map[string]interface{}{"psychologistId": psychologistID}

	cursor, findErr := s.DatabaseUtil.FindMany("treatments", filter)
	if findErr != nil {
		return nil, findErr
	}

	treatments := []*models.GetPsychologistTreatmentsResponse{}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {

		treatment := models.GetPsychologistTreatmentsResponse{}

		decodeErr := cursor.Decode(&treatment)
		if decodeErr != nil {
			return nil, decodeErr
		}

		treatments = append(treatments, &treatment)

	}

	return treatments, nil

}
