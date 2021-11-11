package services

import (
	"github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/identifier"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// GetPsychologistPriceRangeOfferingsService is a service that gets all the price range offerings of a psychologist
type GetPsychologistPriceRangeOfferingsService struct {
	IdentifierUtil identifier.IIdentifierUtil
	OrmUtil        orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s GetPsychologistPriceRangeOfferingsService) Execute(psychologistID string) ([]*models.TreatmentPriceRangeOffering, error) {

	treatmentPriceRangeOfferings := []*models.TreatmentPriceRangeOffering{}

	result := s.OrmUtil.Db().Where("psychologist_id = ?", psychologistID).Find(&treatmentPriceRangeOfferings)
	if result.Error != nil {
		return nil, result.Error
	}

	return treatmentPriceRangeOfferings, nil

}
