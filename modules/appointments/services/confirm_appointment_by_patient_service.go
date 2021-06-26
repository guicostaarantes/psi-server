package services

import (
	"errors"
	"fmt"

	"github.com/guicostaarantes/psi-server/modules/appointments/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// ConfirmAppointmentByPatientService is a service that the patient will use to confirm an appointment
type ConfirmAppointmentByPatientService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s ConfirmAppointmentByPatientService) Execute(id string, patientID string) error {

	appointment := models.Appointment{}

	findErr := s.DatabaseUtil.FindOne("psi_db", "appointments", map[string]interface{}{"id": id, "patientId": patientID}, &appointment)
	if findErr != nil {
		return findErr
	}

	if appointment.ID == "" {
		return errors.New("resource not found")
	}

	if appointment.Status == models.EditedByPatient || appointment.Status == models.ConfirmedByPatient || appointment.Status == models.CanceledByPsychologist {
		return fmt.Errorf("appointment status cannot change from %s to CONFIRMED_BY_PATIENT", string(appointment.Status))
	}

	if appointment.Status == models.EditedByPsychologist || appointment.Status == models.ConfirmedByPsychologist {
		appointment.Status = models.ConfirmedByBoth
	} else {
		appointment.Status = models.ConfirmedByPatient
	}

	appointment.Reason = ""

	writeErr := s.DatabaseUtil.UpdateOne("psi_db", "appointments", map[string]interface{}{"id": id}, appointment)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
