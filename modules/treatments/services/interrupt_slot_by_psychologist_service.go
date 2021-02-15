package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// InterruptSlotByPsychologistService is a service that interrupts a slot, changing its status to interrupted by psychologist
type InterruptSlotByPsychologistService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s InterruptSlotByPsychologistService) Execute(id string, psychologistID string, reason string) error {

	slot := models.Slot{}

	findErr := s.DatabaseUtil.FindOne("psi_db", "slots", map[string]interface{}{"id": id, "psychologistId": psychologistID}, &slot)
	if findErr != nil {
		return findErr
	}

	if slot.ID == "" {
		return errors.New("resource not found")
	}

	if slot.Status != models.Active {
		return fmt.Errorf("slots can only be interrupted if their current status is ACTIVE. current status is %s", string(slot.Status))
	}

	slot.EndDate = time.Now().Unix()
	slot.Status = models.InterruptedByPsychologist
	slot.Reason = reason

	writeErr := s.DatabaseUtil.UpdateOne("psi_db", "slots", map[string]interface{}{"id": id}, slot)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
