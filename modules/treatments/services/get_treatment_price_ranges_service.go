package services

import (
	"context"

	"github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// GetTreatmentPriceRangesService is a service that gets all the treatment price ranges
type GetTreatmentPriceRangesService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s GetTreatmentPriceRangesService) Execute() ([]*models.TreatmentPriceRange, error) {

	priceRanges := []*models.TreatmentPriceRange{}

	cursor, findErr := s.DatabaseUtil.FindMany("treatment_price_ranges", map[string]interface{}{})
	if findErr != nil {
		return nil, findErr
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {

		priceRange := models.TreatmentPriceRange{}

		decodeErr := cursor.Decode(&priceRange)
		if decodeErr != nil {
			return nil, decodeErr
		}

		priceRanges = append(priceRanges, &priceRange)

	}

	return priceRanges, nil

}
