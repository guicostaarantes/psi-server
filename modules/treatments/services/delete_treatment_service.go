package treatments_services

import (
	"errors"

	treatments_models "github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// DeleteTreatmentService is a service that changes data from a treatment
type DeleteTreatmentService struct {
	OrmUtil orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s DeleteTreatmentService) Execute(id string, psychologistID string, priceRangeName string) error {

	treatment := treatments_models.Treatment{}

	result := s.OrmUtil.Db().Where("id = ? AND psychologist_id = ?", id, psychologistID).Limit(1).Find(&treatment)
	if result.Error != nil {
		return result.Error
	}

	if treatment.ID == "" {
		return errors.New("resource not found")
	}

	if treatment.Status != treatments_models.Pending {
		return errors.New("treatments can only be deleted if their status is pending")
	}

	priceRangeOffering := treatments_models.TreatmentPriceRangeOffering{}

	result = s.OrmUtil.Db().Where("psychologist_id = ? AND price_range_name = ?", treatment.PsychologistID, priceRangeName).Limit(1).Find(&priceRangeOffering)
	if result.Error != nil {
		return result.Error
	}

	if priceRangeOffering.ID == "" {
		return errors.New("price range offering not found")
	}

	result = s.OrmUtil.Db().Delete(&treatment)
	if result.Error != nil {
		return result.Error
	}

	result = s.OrmUtil.Db().Delete(&priceRangeOffering)
	if result.Error != nil {
		return result.Error
	}

	return nil

}
