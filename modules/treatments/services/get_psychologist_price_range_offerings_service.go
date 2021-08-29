package services

import (
	"context"

	"github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/identifier"
)

// GetPsychologistPriceRangeOfferingsService is a service that gets all the price range offerings of a psychologist
type GetPsychologistPriceRangeOfferingsService struct {
	DatabaseUtil   database.IDatabaseUtil
	IdentifierUtil identifier.IIdentifierUtil
}

// Execute is the method that runs the business logic of the service
func (s GetPsychologistPriceRangeOfferingsService) Execute(psychologistID string) ([]*models.TreatmentPriceRangeOffering, error) {

	filter := map[string]interface{}{"psychologistId": psychologistID}

	cursor, findErr := s.DatabaseUtil.FindMany("treatment_price_range_offerings", filter)
	if findErr != nil {
		return nil, findErr
	}

	treatments := []*models.TreatmentPriceRangeOffering{}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {

		treatment := models.TreatmentPriceRangeOffering{}

		decodeErr := cursor.Decode(&treatment)
		if decodeErr != nil {
			return nil, decodeErr
		}

		treatments = append(treatments, &treatment)

	}

	return treatments, nil

}
