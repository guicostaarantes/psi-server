package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	characteristic_models "github.com/guicostaarantes/psi-server/modules/characteristics/models"
	"github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// AssignTreatmentService is a service that assigns a patient to a treatment and changes its status to active
type AssignTreatmentService struct {
	DatabaseUtil database.IDatabaseUtil
	OrmUtil      orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s AssignTreatmentService) Execute(id string, priceRangeName string, patientID string) error {

	treatment := models.Treatment{}

	patientInOtherTreatment := models.Treatment{}

	result := s.OrmUtil.Db().Where("patient_id = ? AND status = ?", patientID, models.Active).Limit(1).Find(&patientInOtherTreatment)
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

	if treatment.Status != models.Pending {
		return fmt.Errorf("treatments can only be assigned if their current status is PENDING. current status is %s", string(treatment.Status))
	}

	treatmentPriceRangeOffering := models.TreatmentPriceRangeOffering{}

	findErr := s.DatabaseUtil.FindOne("treatment_price_range_offerings", map[string]interface{}{"psychologistId": treatment.PsychologistID, "priceRangeName": priceRangeName}, &treatmentPriceRangeOffering)
	if findErr != nil {
		return findErr
	}

	if treatmentPriceRangeOffering.ID == "" {
		return errors.New("treatment price range offering not found")
	}

	incomeChar := characteristic_models.CharacteristicChoice{}

	findErr = s.DatabaseUtil.FindOne("characteristic_choices", map[string]interface{}{"profileId": patientID, "characteristicName": "income"}, &incomeChar)
	if findErr != nil {
		return findErr
	}

	if incomeChar.SelectedValue == "" {
		return errors.New("missing income for patient")
	}

	priceRange := models.TreatmentPriceRange{}

	findErr = s.DatabaseUtil.FindOne("treatment_price_ranges", map[string]interface{}{"name": priceRangeName}, &priceRange)
	if findErr != nil {
		return findErr
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
	treatment.Status = models.Active
	treatment.PriceRangeName = priceRangeName

	result = s.OrmUtil.Db().Save(&treatment)
	if result.Error != nil {
		return result.Error
	}

	deleteErr := s.DatabaseUtil.DeleteOne("treatment_price_range_offerings", map[string]interface{}{"id": treatmentPriceRangeOffering.ID})
	if deleteErr != nil {
		return deleteErr
	}

	return nil

}
