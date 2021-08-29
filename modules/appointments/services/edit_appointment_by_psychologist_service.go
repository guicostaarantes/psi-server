package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/guicostaarantes/psi-server/modules/appointments/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// EditAppointmentByPsychologistService is a service that the psychologist will use to edit an appointment
type EditAppointmentByPsychologistService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s EditAppointmentByPsychologistService) Execute(id string, psychologistID string, input models.EditAppointmentByPsychologistInput) error {

	appointment := models.Appointment{}

	findErr := s.DatabaseUtil.FindOne("appointments", map[string]interface{}{"id": id, "psychologistId": psychologistID}, &appointment)
	if findErr != nil {
		return findErr
	}

	if appointment.ID == "" {
		return errors.New("resource not found")
	}

	if appointment.Status == models.CanceledByPatient {
		return fmt.Errorf("appointment status cannot change from %s to EDITED_BY_PSYCHOLOGIST", string(appointment.Status))
	}

	if input.Start < time.Now().Unix() {
		return errors.New("appointment cannot be scheduled to the past")
	}

	if input.Start >= input.End {
		return errors.New("appointment cannot have negative duration")
	}

	appointment.Status = models.EditedByPsychologist
	appointment.Start = input.Start
	appointment.End = input.End
	appointment.PriceRange = input.PriceRange
	appointment.Reason = input.Reason

	writeErr := s.DatabaseUtil.UpdateOne("appointments", map[string]interface{}{"id": id}, appointment)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
