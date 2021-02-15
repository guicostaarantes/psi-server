package services

import (
	"errors"
	"fmt"

	"github.com/guicostaarantes/psi-server/modules/schedule/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// DenyAppointmentService is a service that the psychologist will use to deny an appointment made for him
type DenyAppointmentService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s DenyAppointmentService) Execute(id string, psychologistID string) error {

	appointment := models.Appointment{}

	findErr := s.DatabaseUtil.FindOne("psi_db", "appointments", map[string]interface{}{"id": id, "psychologistId": psychologistID}, &appointment)
	if findErr != nil {
		return findErr
	}

	if appointment.ID == "" {
		return errors.New("resource not found")
	}

	if appointment.Status != models.Proposed {
		return fmt.Errorf("appointments can only be denied if their current status is PROPOSED. current status is %s", string(appointment.Status))
	}

	appointment.Status = models.Denied

	writeErr := s.DatabaseUtil.UpdateOne("psi_db", "appointments", map[string]interface{}{"id": id}, appointment)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
