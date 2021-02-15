package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// AssignSlotService is a service that assigns a patient to a slot and changes its status to active
type AssignSlotService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s AssignSlotService) Execute(id string, patientID string) error {

	slot := models.Slot{}

	patientInOtherSlot := models.Slot{}

	findErr := s.DatabaseUtil.FindOne("psi_db", "slots", map[string]interface{}{"patientId": patientID, "status": string(models.Active)}, &patientInOtherSlot)
	if findErr != nil {
		return findErr
	}

	if patientInOtherSlot.ID != "" {
		return errors.New("patient is already in an active slot")
	}

	findErr = s.DatabaseUtil.FindOne("psi_db", "slots", map[string]interface{}{"id": id}, &slot)
	if findErr != nil {
		return findErr
	}

	if slot.ID == "" {
		return errors.New("resource not found")
	}

	if slot.Status != models.Pending {
		return fmt.Errorf("slots can only be assigned if their current status is PENDING. current status is %s", string(slot.Status))
	}

	slot.PatientID = patientID
	slot.StartDate = time.Now().Unix()
	slot.Status = models.Active

	writeErr := s.DatabaseUtil.UpdateOne("psi_db", "slots", map[string]interface{}{"id": id}, slot)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
