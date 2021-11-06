package services

import (
	"github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// GetTreatmentPriceRangesService is a service that gets all the treatment price ranges
type GetTreatmentPriceRangesService struct {
	OrmUtil orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s GetTreatmentPriceRangesService) Execute() ([]*models.TreatmentPriceRange, error) {

	priceRanges := []*models.TreatmentPriceRange{}

	result := s.OrmUtil.Db().Find(&priceRanges)
	if result.Error != nil {
		return nil, result.Error
	}

	return priceRanges, nil

}
