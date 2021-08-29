package services

import (
	"github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/identifier"
)

// CreateTreatmentService is a service that creates a new treatment for a psychologist
type CreateTreatmentService struct {
	DatabaseUtil            database.IDatabaseUtil
	IdentifierUtil          identifier.IIdentifierUtil
	ScheduleIntervalSeconds int64
}

// Execute is the method that runs the business logic of the service
func (s CreateTreatmentService) Execute(psychologistID string, input models.CreateTreatmentInput) error {

	checkTreatmentCollisionService := CheckTreatmentCollisionService{
		DatabaseUtil:            s.DatabaseUtil,
		ScheduleIntervalSeconds: s.ScheduleIntervalSeconds,
	}

	checkErr := checkTreatmentCollisionService.Execute(psychologistID, input.Frequency, input.Phase, input.Duration, "")
	if checkErr != nil {
		return checkErr
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
		PriceRange:     input.PriceRange,
	}

	writeErr = s.DatabaseUtil.InsertOne("treatment_price_offerings", treatmentPriceOffering)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
