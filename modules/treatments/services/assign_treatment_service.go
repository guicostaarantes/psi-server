package treatments_services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	characteristic_models "github.com/guicostaarantes/psi-server/modules/characteristics/models"
	cooldowns_models "github.com/guicostaarantes/psi-server/modules/cooldowns/models"
	cooldowns_services "github.com/guicostaarantes/psi-server/modules/cooldowns/services"
	treatments_models "github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// AssignTreatmentService is a service that assigns a patient to a treatment and changes its status to active
type AssignTreatmentService struct {
	OrmUtil            orm.IOrmUtil
	GetCooldownService *cooldowns_services.GetCooldownService
}

// Execute is the method that runs the business logic of the service
func (s AssignTreatmentService) Execute(id string, priceRangeName string, patientID string) error {

	cooldown, getErr := s.GetCooldownService.Execute(patientID, cooldowns_models.Patient, cooldowns_models.TreatmentInterrupted)
	if getErr != nil {
		return getErr
	}

	if cooldown != nil {
		return fmt.Errorf("assign treatment is blocked for this user until %d", cooldown.ValidUntil)
	}

	treatment := treatments_models.Treatment{}

	patientInOtherTreatment := treatments_models.Treatment{}

	result := s.OrmUtil.Db().Where("patient_id = ? AND status = ?", patientID, treatments_models.Active).Limit(1).Find(&patientInOtherTreatment)
	if result.Error != nil {
		return result.Error
	}

	if patientInOtherTreatment.ID != "" {
		return errors.New("patient is already in an active treatment")
	}

	result = s.OrmUtil.Db().Where("id = ?", id).Limit(1).Find(&treatment)
	if result.Error != nil {
		return result.Error
	}

	if treatment.ID == "" {
		return errors.New("resource not found")
	}

	if treatment.Status != treatments_models.Pending {
		return fmt.Errorf("treatments can only be assigned if their current status is PENDING. current status is %s", string(treatment.Status))
	}

	treatmentPriceRangeOffering := treatments_models.TreatmentPriceRangeOffering{}

	result = s.OrmUtil.Db().Where("psychologist_id = ? AND price_range_name = ?", treatment.PsychologistID, priceRangeName).Limit(1).Find(&treatmentPriceRangeOffering)
	if result.Error != nil {
		return result.Error
	}

	if treatmentPriceRangeOffering.ID == "" {
		return errors.New("treatment price range offering not found")
	}

	incomeChar := characteristic_models.CharacteristicChoice{}

	result = s.OrmUtil.Db().Where("profile_id = ? AND characteristic_name = ?", patientID, "income").Limit(1).Find(&incomeChar)
	if result.Error != nil {
		return result.Error
	}

	if incomeChar.SelectedValue == "" {
		return errors.New("missing income for patient")
	}

	priceRange := treatments_models.TreatmentPriceRange{}

	result = s.OrmUtil.Db().Where("name = ?", priceRangeName).Limit(1).Find(&priceRange)
	if result.Error != nil {
		return result.Error
	}

	if priceRange.EligibleFor == "" {
		return errors.New("missing price range eligibility parameters")
	}

	isEligible := false

	eligibleParameters := strings.Split(priceRange.EligibleFor, ",")
	for _, v := range eligibleParameters {
		if v == incomeChar.SelectedValue {
			isEligible = true
		}
	}

	if !isEligible {
		return errors.New("patient is not eligible for this price range")
	}

	treatment.PatientID = patientID
	treatment.StartDate = time.Now().Unix()
	treatment.Status = treatments_models.Active
	treatment.PriceRangeName = priceRangeName

	result = s.OrmUtil.Db().Save(&treatment)
	if result.Error != nil {
		return result.Error
	}

	result = s.OrmUtil.Db().Delete(&treatmentPriceRangeOffering)
	if result.Error != nil {
		return result.Error
	}

	return nil

}
