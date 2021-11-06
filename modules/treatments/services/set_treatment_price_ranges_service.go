package services

import (
	"github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// SetTreatmentPriceRangesService is a service that sets all possible patient characteristics
type SetTreatmentPriceRangesService struct {
	OrmUtil orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s SetTreatmentPriceRangesService) Execute(input []*models.TreatmentPriceRange) error {

	result := s.OrmUtil.Db().Where("1 = 1").Delete(&models.TreatmentPriceRange{})
	if result.Error != nil {
		return result.Error
	}

	result = s.OrmUtil.Db().Create(&input)
	if result.Error != nil {
		return result.Error
	}

	return nil

}
