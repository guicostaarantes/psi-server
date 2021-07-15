package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/guicostaarantes/psi-server/modules/appointments/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// EditAppointmentByPatientService is a service that the patient will use to edit an appointment
type EditAppointmentByPatientService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s EditAppointmentByPatientService) Execute(id string, patientID string, input models.EditAppointmentByPatientInput) error {

	appointment := models.Appointment{}

	findErr := s.DatabaseUtil.FindOne("psi_db", "appointments", map[string]interface{}{"id": id, "patientId": patientID}, &appointment)
	if findErr != nil {
		return findErr
	}

	if appointment.ID == "" {
		return errors.New("resource not found")
	}

	if appointment.Status == models.CanceledByPsychologist {
		return fmt.Errorf("appointment status cannot change from %s to EDITED_BY_PATIENT", string(appointment.Status))
	}

	if input.Start < time.Now().Unix() {
		return errors.New("appointment cannot be scheduled to the past")
	}

	appointment.Status = models.EditedByPatient
	appointment.End += input.Start - appointment.Start
	appointment.Start = input.Start
	appointment.Reason = input.Reason

	writeErr := s.DatabaseUtil.UpdateOne("psi_db", "appointments", map[string]interface{}{"id": id}, appointment)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
