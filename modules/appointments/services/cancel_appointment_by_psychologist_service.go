package services

import (
	"errors"
	"fmt"

	"github.com/guicostaarantes/psi-server/modules/appointments/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// CancelAppointmentByPsychologistService is a service that the psychologist will use to cancel an appointment
type CancelAppointmentByPsychologistService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s CancelAppointmentByPsychologistService) Execute(id string, psychologistID string, reason string) error {

	appointment := models.Appointment{}

	findErr := s.DatabaseUtil.FindOne("psi_db", "appointments", map[string]interface{}{"id": id, "psychologistId": psychologistID}, &appointment)
	if findErr != nil {
		return findErr
	}

	if appointment.ID == "" {
		return errors.New("resource not found")
	}

	if appointment.Status == models.CanceledByPatient || appointment.Status == models.CanceledByPsychologist {
		return fmt.Errorf("appointment status cannot change from %s to CANCELED_BY_PSYCHOLOGIST", string(appointment.Status))
	}

	appointment.Status = models.CanceledByPsychologist
	appointment.Reason = reason

	writeErr := s.DatabaseUtil.UpdateOne("psi_db", "appointments", map[string]interface{}{"id": id}, appointment)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
