package services

import (
	"github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/identifier"
)

// CreateTreatmentService is a service that creates a new treatment for a psychologist
type CreateTreatmentService struct {
	DatabaseUtil   database.IDatabaseUtil
	IdentifierUtil identifier.IIdentifierUtil
}

// Execute is the method that runs the business logic of the service
func (s CreateTreatmentService) Execute(psychologistID string, input models.CreateTreatmentInput) error {

	_, treatmentID, treatmentIDErr := s.IdentifierUtil.GenerateIdentifier()
	if treatmentIDErr != nil {
		return treatmentIDErr
	}

	treatment := models.Treatment{
		ID:             treatmentID,
		PsychologistID: psychologistID,
		Duration:       input.Duration,
		Price:          input.Price,
		Interval:       input.Interval,
		Status:         models.Pending,
	}

	writeErr := s.DatabaseUtil.InsertOne("psi_db", "treatments", treatment)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
