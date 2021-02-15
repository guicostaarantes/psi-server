package services

import (
	"github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/identifier"
)

// CreateSlotService is a service that creates a new slot for a psychologist
type CreateSlotService struct {
	DatabaseUtil   database.IDatabaseUtil
	IdentifierUtil identifier.IIdentifierUtil
}

// Execute is the method that runs the business logic of the service
func (s CreateSlotService) Execute(psychologistID string, input models.CreateSlotInput) error {

	_, slotID, slotIDErr := s.IdentifierUtil.GenerateIdentifier()
	if slotIDErr != nil {
		return slotIDErr
	}

	slot := models.Slot{
		ID:             slotID,
		PsychologistID: psychologistID,
		Duration:       input.Duration,
		Price:          input.Price,
		Interval:       input.Interval,
		Status:         models.Pending,
	}

	writeErr := s.DatabaseUtil.InsertOne("psi_db", "slots", slot)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
