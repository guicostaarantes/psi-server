package services

import (
	"github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// GetTreatmentPriceRangeByNameService is a service that gets a treatment based on its id
type GetTreatmentPriceRangeByNameService struct {
	OrmUtil orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s GetTreatmentPriceRangeByNameService) Execute(name string) (*models.TreatmentPriceRange, error) {

	priceRange := &models.TreatmentPriceRange{}

	result := s.OrmUtil.Db().Where("name = ?", name).Limit(1).Find(&priceRange)
	if result.Error != nil {
		return nil, result.Error
	}

	return priceRange, nil

}
