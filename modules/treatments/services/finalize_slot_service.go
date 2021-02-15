package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// FinalizeSlotService is a service that changes the status of a slot to finalized
type FinalizeSlotService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s FinalizeSlotService) Execute(id string, psychologistID string) error {

	slot := models.Slot{}

	findErr := s.DatabaseUtil.FindOne("psi_db", "slots", map[string]interface{}{"id": id, "psychologistId": psychologistID}, &slot)
	if findErr != nil {
		return findErr
	}

	if slot.ID == "" {
		return errors.New("resource not found")
	}

	if slot.Status != models.Active {
		return fmt.Errorf("slots can only be finalized if their current status is ACTIVE. current status is %s", string(slot.Status))
	}

	slot.EndDate = time.Now().Unix()
	slot.Status = models.Finalized

	writeErr := s.DatabaseUtil.UpdateOne("psi_db", "slots", map[string]interface{}{"id": id}, slot)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
