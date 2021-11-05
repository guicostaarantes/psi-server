package services

import (
	"errors"

	"github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// DeleteTreatmentService is a service that changes data from a treatment
type DeleteTreatmentService struct {
	DatabaseUtil database.IDatabaseUtil
	OrmUtil      orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s DeleteTreatmentService) Execute(id string, psychologistID string, priceRangeName string) error {

	treatment := models.Treatment{}

	result := s.OrmUtil.Db().Where("id = ? AND psychologist_id = ?", id, psychologistID).Limit(1).Find(&treatment)
	if result.Error != nil {
		return result.Error
	}

	if treatment.ID == "" {
		return errors.New("resource not found")
	}

	if treatment.Status != models.Pending {
		return errors.New("treatments can only be deleted if their status is pending")
	}

	priceRangeOffering := models.TreatmentPriceRangeOffering{}

	findErr := s.DatabaseUtil.FindOne("treatment_price_range_offerings", map[string]interface{}{"priceRangeName": priceRangeName, "psychologistId": psychologistID}, &priceRangeOffering)
	if findErr != nil {
		return findErr
	}

	if priceRangeOffering.ID == "" {
		return errors.New("price range offering not found")
	}

	result = s.OrmUtil.Db().Delete(&treatment)
	if result.Error != nil {
		return result.Error
	}

	writeErr := s.DatabaseUtil.DeleteOne("treatment_price_range_offerings", map[string]interface{}{"id": priceRangeOffering.ID})
	if writeErr != nil {
		return writeErr
	}

	return nil

}
