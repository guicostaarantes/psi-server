package services

import (
	"github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// SetTreatmentPriceRangesService is a service that sets all possible patient characteristics
type SetTreatmentPriceRangesService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s SetTreatmentPriceRangesService) Execute(input []*models.TreatmentPriceRange) error {

	newPriceRanges := []interface{}{}

	for _, pr := range input {
		priceRange := models.TreatmentPriceRange{
			Name:         pr.Name,
			MinimumPrice: pr.MinimumPrice,
			MaximumPrice: pr.MaximumPrice,
			EligibleFor:  pr.EligibleFor,
		}

		newPriceRanges = append(newPriceRanges, priceRange)
	}

	deleteErr := s.DatabaseUtil.DeleteMany("treatment_price_ranges", map[string]interface{}{})
	if deleteErr != nil {
		return deleteErr
	}

	writeErr := s.DatabaseUtil.InsertMany("treatment_price_ranges", newPriceRanges)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
