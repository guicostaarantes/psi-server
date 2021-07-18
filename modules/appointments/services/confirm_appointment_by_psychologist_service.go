package services

import (
	"errors"
	"fmt"

	"github.com/guicostaarantes/psi-server/modules/appointments/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// ConfirmAppointmentByPsychologistService is a service that the psychologist will use to confirm an appointment
type ConfirmAppointmentByPsychologistService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s ConfirmAppointmentByPsychologistService) Execute(id string, psychologistID string) error {

	appointment := models.Appointment{}

	findErr := s.DatabaseUtil.FindOne("appointments", map[string]interface{}{"id": id, "psychologistId": psychologistID}, &appointment)
	if findErr != nil {
		return findErr
	}

	if appointment.ID == "" {
		return errors.New("resource not found")
	}

	if appointment.Status == models.EditedByPsychologist || appointment.Status == models.ConfirmedByPsychologist || appointment.Status == models.CanceledByPatient {
		return fmt.Errorf("appointment status cannot change from %s to CONFIRMED_BY_PSYCHOLOGIST", string(appointment.Status))
	}

	if appointment.Status == models.EditedByPatient || appointment.Status == models.ConfirmedByPatient {
		appointment.Status = models.ConfirmedByBoth
	} else {
		appointment.Status = models.ConfirmedByPsychologist
	}

	appointment.Reason = ""

	writeErr := s.DatabaseUtil.UpdateOne("appointments", map[string]interface{}{"id": id}, appointment)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
