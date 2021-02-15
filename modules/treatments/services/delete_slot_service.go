package services

import (
	"errors"

	"github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// DeleteSlotService is a service that changes data from a slot
type DeleteSlotService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s DeleteSlotService) Execute(id string, psychologistID string) error {

	slot := models.Slot{}

	findErr := s.DatabaseUtil.FindOne("psi_db", "slots", map[string]interface{}{"id": id, "psychologistId": psychologistID}, &slot)
	if findErr != nil {
		return findErr
	}

	if slot.ID == "" {
		return errors.New("resource not found")
	}

	if slot.Status != models.Pending {
		return errors.New("slots can only be deleted if they their status is pending")
	}

	writeErr := s.DatabaseUtil.DeleteOne("psi_db", "slots", map[string]interface{}{"id": id})
	if writeErr != nil {
		return writeErr
	}

	return nil

}
