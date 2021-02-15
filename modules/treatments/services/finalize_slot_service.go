package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// FinalizeTreatmentService is a service that changes the status of a treatment to finalized
type FinalizeTreatmentService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s FinalizeTreatmentService) Execute(id string, psychologistID string) error {

	treatment := models.Treatment{}

	findErr := s.DatabaseUtil.FindOne("psi_db", "treatments", map[string]interface{}{"id": id, "psychologistId": psychologistID}, &treatment)
	if findErr != nil {
		return findErr
	}

	if treatment.ID == "" {
		return errors.New("resource not found")
	}

	if treatment.Status != models.Active {
		return fmt.Errorf("treatments can only be finalized if their current status is ACTIVE. current status is %s", string(treatment.Status))
	}

	treatment.EndDate = time.Now().Unix()
	treatment.Status = models.Finalized

	writeErr := s.DatabaseUtil.UpdateOne("psi_db", "treatments", map[string]interface{}{"id": id}, treatment)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
