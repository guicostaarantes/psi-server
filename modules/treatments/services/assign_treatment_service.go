package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	characteristic_models "github.com/guicostaarantes/psi-server/modules/characteristics/models"
	"github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// AssignTreatmentService is a service that assigns a patient to a treatment and changes its status to active
type AssignTreatmentService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s AssignTreatmentService) Execute(id string, priceRangeName string, patientID string) error {

	treatment := models.Treatment{}

	patientInOtherTreatment := models.Treatment{}

	findErr := s.DatabaseUtil.FindOne("treatments", map[string]interface{}{"patientId": patientID, "status": string(models.Active)}, &patientInOtherTreatment)
	if findErr != nil {
		return findErr
	}

	if patientInOtherTreatment.ID != "" {
		return errors.New("patient is already in an active treatment")
	}

	findErr = s.DatabaseUtil.FindOne("treatments", map[string]interface{}{"id": id}, &treatment)
	if findErr != nil {
		return findErr
	}

	if treatment.ID == "" {
		return errors.New("resource not found")
	}

	if treatment.Status != models.Pending {
		return fmt.Errorf("treatments can only be assigned if their current status is PENDING. current status is %s", string(treatment.Status))
	}

	treatmentPriceRangeOffering := models.TreatmentPriceRangeOffering{}

	findErr = s.DatabaseUtil.FindOne("treatment_price_range_offerings", map[string]interface{}{"psychologistId": treatment.PsychologistID, "priceRange": priceRangeName}, &treatmentPriceRangeOffering)
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
	treatment.PriceRange = priceRangeName

	writeErr := s.DatabaseUtil.UpdateOne("treatments", map[string]interface{}{"id": id}, treatment)
	if writeErr != nil {
		return writeErr
	}

	deleteErr := s.DatabaseUtil.DeleteOne("treatment_price_range_offerings", map[string]interface{}{"id": treatmentPriceRangeOffering.ID})
	if deleteErr != nil {
		return deleteErr
	}

	return nil

}
