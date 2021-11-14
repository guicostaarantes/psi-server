package treatments_services

import (
	"errors"

	treatments_models "github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/identifier"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// CreateTreatmentService is a service that creates a new treatment for a psychologist
type CreateTreatmentService struct {
	IdentifierUtil                 identifier.IIdentifierUtil
	OrmUtil                        orm.IOrmUtil
	CheckTreatmentCollisionService *CheckTreatmentCollisionService
}

// Execute is the method that runs the business logic of the service
func (s CreateTreatmentService) Execute(psychologistID string, input treatments_models.CreateTreatmentInput) error {

	checkErr := s.CheckTreatmentCollisionService.Execute(psychologistID, input.Frequency, input.Phase, input.Duration, "")
	if checkErr != nil {
		return checkErr
	}

	priceRange := treatments_models.TreatmentPriceRange{}

	result := s.OrmUtil.Db().Where("name = ?", input.PriceRangeName).Limit(1).Find(&priceRange)
	if result.Error != nil {
		return result.Error
	}

	if priceRange.Name == "" {
		return errors.New("price range name not found")
	}

	_, treatmentID, treatmentIDErr := s.IdentifierUtil.GenerateIdentifier()
	if treatmentIDErr != nil {
		return treatmentIDErr
	}

	treatment := treatments_models.Treatment{
		ID:             treatmentID,
		PsychologistID: psychologistID,
		Frequency:      input.Frequency,
		Phase:          input.Phase,
		Duration:       input.Duration,
		Status:         treatments_models.Pending,
	}

	result = s.OrmUtil.Db().Create(&treatment)
	if result.Error != nil {
		return result.Error
	}

	_, treatmentPriceOfferingID, treatmentPriceOfferingIDErr := s.IdentifierUtil.GenerateIdentifier()
	if treatmentPriceOfferingIDErr != nil {
		return treatmentPriceOfferingIDErr
	}

	treatmentPriceOffering := treatments_models.TreatmentPriceRangeOffering{
		ID:             treatmentPriceOfferingID,
		PsychologistID: psychologistID,
		PriceRangeName: input.PriceRangeName,
	}

	result = s.OrmUtil.Db().Create(&treatmentPriceOffering)
	if result.Error != nil {
		return result.Error
	}

	return nil

}
