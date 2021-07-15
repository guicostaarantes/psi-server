package services

import (
	"context"

	"github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/identifier"
)

// GetPsychologistPendingTreatmentsService is a service that gets all the pending treatments of a psychologist
type GetPsychologistPendingTreatmentsService struct {
	DatabaseUtil   database.IDatabaseUtil
	IdentifierUtil identifier.IIdentifierUtil
}

// Execute is the method that runs the business logic of the service
func (s GetPsychologistPendingTreatmentsService) Execute(psychologistID string) ([]*models.GetPsychologistTreatmentsResponse, error) {

	filter := map[string]interface{}{"psychologistId": psychologistID}

	cursor, findErr := s.DatabaseUtil.FindMany("psi_db", "treatments", filter)
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

		if treatment.Status == models.Pending {
			treatments = append(treatments, &treatment)
		}

	}

	return treatments, nil

}
