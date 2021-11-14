package appointments_services

import (
	"errors"
	"fmt"

	appointments_models "github.com/guicostaarantes/psi-server/modules/appointments/models"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// ConfirmAppointmentByPsychologistService is a service that the psychologist will use to confirm an appointment
type ConfirmAppointmentByPsychologistService struct {
	OrmUtil orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s ConfirmAppointmentByPsychologistService) Execute(id string, psychologistID string) error {

	appointment := appointments_models.Appointment{}

	result := s.OrmUtil.Db().Where("id = ? AND psychologist_id = ?", id, psychologistID).Limit(1).Find(&appointment)
	if result.Error != nil {
		return result.Error
	}

	if appointment.ID == "" {
		return errors.New("resource not found")
	}

	if appointment.Status == appointments_models.EditedByPsychologist || appointment.Status == appointments_models.ConfirmedByPsychologist || appointment.Status == appointments_models.CanceledByPatient {
		return fmt.Errorf("appointment status cannot change from %s to CONFIRMED_BY_PSYCHOLOGIST", string(appointment.Status))
	}

	if appointment.Status == appointments_models.EditedByPatient || appointment.Status == appointments_models.ConfirmedByPatient {
		appointment.Status = appointments_models.ConfirmedByBoth
	} else {
		appointment.Status = appointments_models.ConfirmedByPsychologist
	}

	appointment.Reason = ""

	result = s.OrmUtil.Db().Save(&appointment)
	if result.Error != nil {
		return result.Error
	}

	return nil

}
