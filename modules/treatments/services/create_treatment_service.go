package services

import (
	"errors"

	"github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/identifier"
)

// CreateTreatmentService is a service that creates a new treatment for a psychologist
type CreateTreatmentService struct {
	DatabaseUtil                   database.IDatabaseUtil
	IdentifierUtil                 identifier.IIdentifierUtil
	CheckTreatmentCollisionService *CheckTreatmentCollisionService
}

// Execute is the method that runs the business logic of the service
func (s CreateTreatmentService) Execute(psychologistID string, input models.CreateTreatmentInput) error {

	checkErr := s.CheckTreatmentCollisionService.Execute(psychologistID, input.Frequency, input.Phase, input.Duration, "")
	if checkErr != nil {
		return checkErr
	}

	priceRange := models.TreatmentPriceRange{}

	findErr := s.DatabaseUtil.FindOne("treatment_price_ranges", map[string]interface{}{"name": input.PriceRangeName}, &priceRange)
	if findErr != nil {
		return findErr
	}

	if priceRange.Name == "" {
		return errors.New("price range name not found")
	}

	_, treatmentID, treatmentIDErr := s.IdentifierUtil.GenerateIdentifier()
	if treatmentIDErr != nil {
		return treatmentIDErr
	}

	treatment := models.Treatment{
		ID:             treatmentID,
		PsychologistID: psychologistID,
		Frequency:      input.Frequency,
		Phase:          input.Phase,
		Duration:       input.Duration,
		Status:         models.Pending,
	}

	writeErr := s.DatabaseUtil.InsertOne("treatments", treatment)
	if writeErr != nil {
		return writeErr
	}

	_, treatmentPriceOfferingID, treatmentPriceOfferingIDErr := s.IdentifierUtil.GenerateIdentifier()
	if treatmentPriceOfferingIDErr != nil {
		return treatmentPriceOfferingIDErr
	}

	treatmentPriceOffering := models.TreatmentPriceRangeOffering{
		ID:             treatmentPriceOfferingID,
		PsychologistID: psychologistID,
		PriceRangeName: input.PriceRangeName,
	}

	writeErr = s.DatabaseUtil.InsertOne("treatment_price_range_offerings", treatmentPriceOffering)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
