package services

import (
	"errors"

	"github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// UpdateTreatmentService is a service that changes data from a treatment
type UpdateTreatmentService struct {
	OrmUtil                        orm.IOrmUtil
	CheckTreatmentCollisionService *CheckTreatmentCollisionService
}

// Execute is the method that runs the business logic of the service
func (s UpdateTreatmentService) Execute(id string, psychologistID string, input models.UpdateTreatmentInput) error {

	treatment := models.Treatment{}

	result := s.OrmUtil.Db().Where("id = ? AND psychologist_id = ?", id, psychologistID).Limit(1).Find(&treatment)
	if result.Error != nil {
		return result.Error
	}

	if treatment.ID == "" {
		return errors.New("resource not found")
	}

	if input.PriceRangeName != "" && treatment.Status == models.Pending {
		return errors.New("pending treatments are not allowed to have a price range")
	}

	if input.PriceRangeName == "" && treatment.Status != models.Pending {
		return errors.New("non-pending treatments must have a price range")
	}

	checkErr := s.CheckTreatmentCollisionService.Execute(psychologistID, input.Frequency, input.Phase, input.Duration, id)
	if checkErr != nil {
		return checkErr
	}

	treatment.Frequency = input.Frequency
	treatment.Phase = input.Phase
	treatment.Duration = input.Duration
	treatment.PriceRangeName = input.PriceRangeName

	result = s.OrmUtil.Db().Save(&treatment)
	if result.Error != nil {
		return result.Error
	}

	return nil

}
