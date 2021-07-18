package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	appointments_models "github.com/guicostaarantes/psi-server/modules/appointments/models"
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

	findErr := s.DatabaseUtil.FindOne("treatments", map[string]interface{}{"id": id, "psychologistId": psychologistID}, &treatment)
	if findErr != nil {
		return findErr
	}

	if treatment.ID == "" {
		return errors.New("resource not found")
	}

	if treatment.Status != models.Active {
		return fmt.Errorf("treatments can only be interrupted if their current status is ACTIVE. current status is %s", string(treatment.Status))
	}

	cursor, findErr := s.DatabaseUtil.FindMany("appointments", map[string]interface{}{"treatmentId": id})
	if findErr != nil {
		return findErr
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		appointment := appointments_models.Appointment{}

		decodeErr := cursor.Decode(&appointment)
		if decodeErr != nil {
			return decodeErr
		}

		if appointment.Start > time.Now().Unix() && appointment.Status != appointments_models.CanceledByPatient {
			appointment.Status = appointments_models.TreatmentInterruptedByPsychologist
			appointment.Reason = reason

			writeErr := s.DatabaseUtil.UpdateOne("appointments", map[string]interface{}{"id": appointment.ID}, appointment)
			if writeErr != nil {
				return writeErr
			}
		}
	}

	treatment.EndDate = time.Now().Unix()
	treatment.Status = models.InterruptedByPsychologist
	treatment.Reason = reason

	writeErr := s.DatabaseUtil.UpdateOne("treatments", map[string]interface{}{"id": id}, treatment)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
