package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	schedule_models "github.com/guicostaarantes/psi-server/modules/schedule/models"
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

	cursor, findErr := s.DatabaseUtil.FindMany("psi_db", "appointments", map[string]interface{}{"treatmentId": id})
	if findErr != nil {
		return findErr
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		appointment := schedule_models.Appointment{}

		decodeErr := cursor.Decode(&appointment)
		if decodeErr != nil {
			return decodeErr
		}

		if appointment.Start > time.Now().Unix() && (appointment.Status == schedule_models.Proposed || appointment.Status == schedule_models.Confirmed) {
			appointment.Status = schedule_models.CanceledByPsychologist
			appointment.Reason = "Finished treatment"

			writeErr := s.DatabaseUtil.UpdateOne("psi_db", "appointments", map[string]interface{}{"id": appointment.ID}, appointment)
			if writeErr != nil {
				return writeErr
			}
		}
	}

	treatment.EndDate = time.Now().Unix()
	treatment.Status = models.Finalized

	writeErr := s.DatabaseUtil.UpdateOne("psi_db", "treatments", map[string]interface{}{"id": id}, treatment)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
