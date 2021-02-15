package services

import (
	"errors"
	"fmt"

	"github.com/guicostaarantes/psi-server/modules/schedule/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// ConfirmAppointmentService is a service that the psychologist will use to deny an appointment made for him
type ConfirmAppointmentService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s ConfirmAppointmentService) Execute(id string, psychologistID string) error {

	appointment := models.Appointment{}

	findErr := s.DatabaseUtil.FindOne("psi_db", "appointments", map[string]interface{}{"id": id, "psychologistId": psychologistID}, &appointment)
	if findErr != nil {
		return findErr
	}

	if appointment.ID == "" {
		return errors.New("resource not found")
	}

	if appointment.Status != models.Proposed {
		return fmt.Errorf("appointments can only be confirmed if their current status is PROPOSED. current status is %s", string(appointment.Status))
	}

	appointment.Status = models.Confirmed

	writeErr := s.DatabaseUtil.UpdateOne("psi_db", "appointments", map[string]interface{}{"id": id}, appointment)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
