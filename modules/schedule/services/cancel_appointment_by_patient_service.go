package services

import (
	"errors"
	"fmt"

	"github.com/guicostaarantes/psi-server/modules/schedule/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// CancelAppointmentByPatientService is a service that the patient will use to cancel an appointment made for him
type CancelAppointmentByPatientService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s CancelAppointmentByPatientService) Execute(id string, patientID string, reason string) error {

	appointment := models.Appointment{}

	findErr := s.DatabaseUtil.FindOne("psi_db", "appointments", map[string]interface{}{"id": id, "patientId": patientID}, &appointment)
	if findErr != nil {
		return findErr
	}

	if appointment.ID == "" {
		return errors.New("resource not found")
	}

	if appointment.Status != models.Confirmed {
		return fmt.Errorf("appointments can only be canceled if their current status is CONFIRMED. current status is %s", string(appointment.Status))
	}

	appointment.Status = models.CanceledByPatient
	appointment.Reason = reason

	writeErr := s.DatabaseUtil.UpdateOne("psi_db", "appointments", map[string]interface{}{"id": id}, appointment)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
