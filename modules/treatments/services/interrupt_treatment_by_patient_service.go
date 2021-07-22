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

// InterruptTreatmentByPatientService is a service that interrupts a treatment, changing its status to interrupted by patient
type InterruptTreatmentByPatientService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s InterruptTreatmentByPatientService) Execute(id string, patientID string, reason string) error {

	treatment := models.Treatment{}

	findErr := s.DatabaseUtil.FindOne("treatments", map[string]interface{}{"id": id, "patientId": patientID}, &treatment)
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

		if appointment.Start > time.Now().Unix() && appointment.Status != appointments_models.CanceledByPsychologist {
			appointment.Status = appointments_models.TreatmentInterruptedByPatient
			appointment.Reason = reason

			writeErr := s.DatabaseUtil.UpdateOne("appointments", map[string]interface{}{"id": appointment.ID}, appointment)
			if writeErr != nil {
				return writeErr
			}
		}
	}

	treatment.EndDate = time.Now().Unix()
	treatment.Status = models.InterruptedByPatient
	treatment.Reason = reason

	writeErr := s.DatabaseUtil.UpdateOne("treatments", map[string]interface{}{"id": id}, treatment)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
