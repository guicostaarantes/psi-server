package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// InterruptTreatmentByPsychologistService is a service that interrupts a treatment, changing its status to interrupted by psychologist
type InterruptTreatmentByPsychologistService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s InterruptTreatmentByPsychologistService) Execute(id string, psychologistID string, reason string) error {

	treatment := models.Treatment{}

	findErr := s.DatabaseUtil.FindOne("psi_db", "treatments", map[string]interface{}{"id": id, "psychologistId": psychologistID}, &treatment)
	if findErr != nil {
		return findErr
	}

	if treatment.ID == "" {
		return errors.New("resource not found")
	}

	if treatment.Status != models.Active {
		return fmt.Errorf("treatments can only be interrupted if their current status is ACTIVE. current status is %s", string(treatment.Status))
	}

	treatment.EndDate = time.Now().Unix()
	treatment.Status = models.InterruptedByPsychologist
	treatment.Reason = reason

	writeErr := s.DatabaseUtil.UpdateOne("psi_db", "treatments", map[string]interface{}{"id": id}, treatment)
	if writeErr != nil {
		return writeErr
	}

	return nil

}