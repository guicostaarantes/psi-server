package services

import (
	"github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// GetTreatmentPriceRangeByNameService is a service that gets a treatment based on its id
type GetTreatmentPriceRangeByNameService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s GetTreatmentPriceRangeByNameService) Execute(name string) (*models.TreatmentPriceRange, error) {

	priceRange := &models.TreatmentPriceRange{}

	findErr := s.DatabaseUtil.FindOne("treatment_price_ranges", map[string]interface{}{"name": name}, &priceRange)
	if findErr != nil {
		return nil, findErr
	}

	return priceRange, nil

}
