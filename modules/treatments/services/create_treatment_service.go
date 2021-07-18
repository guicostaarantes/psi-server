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

	checkTreatmentCollisionService := CheckTreatmentCollisionService{
		DatabaseUtil: s.DatabaseUtil,
	}

	checkErr := checkTreatmentCollisionService.Execute(psychologistID, input.WeeklyStart, input.Duration, "")
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
		WeeklyStart:    input.WeeklyStart,
		Duration:       input.Duration,
		Price:          input.Price,
		Status:         models.Pending,
	}

	writeErr := s.DatabaseUtil.InsertOne("treatments", treatment)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
