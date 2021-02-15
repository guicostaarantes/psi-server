package services

import (
	"errors"

	"github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// UpdateSlotService is a service that changes data from a slot
type UpdateSlotService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s UpdateSlotService) Execute(id string, psychologistID string, input models.UpdateSlotInput) error {

	slot := models.Slot{}

	findErr := s.DatabaseUtil.FindOne("psi_db", "slots", map[string]interface{}{"id": id, "psychologistId": psychologistID}, &slot)
	if findErr != nil {
		return findErr
	}

	if slot.ID == "" {
		return errors.New("resource not found")
	}

	slot.Duration = input.Duration
	slot.Price = input.Price
	slot.Interval = input.Interval

	writeErr := s.DatabaseUtil.UpdateOne("psi_db", "slots", map[string]interface{}{"id": id}, slot)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
