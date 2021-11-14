package appointments_services

import (
	"errors"
	"fmt"
	"time"

	appointments_models "github.com/guicostaarantes/psi-server/modules/appointments/models"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// EditAppointmentByPsychologistService is a service that the psychologist will use to edit an appointment
type EditAppointmentByPsychologistService struct {
	OrmUtil orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s EditAppointmentByPsychologistService) Execute(id string, psychologistID string, input appointments_models.EditAppointmentByPsychologistInput) error {

	appointment := appointments_models.Appointment{}

	result := s.OrmUtil.Db().Where("id = ? AND psychologist_id = ?", id, psychologistID).Limit(1).Find(&appointment)
	if result.Error != nil {
		return result.Error
	}

	if appointment.ID == "" {
		return errors.New("resource not found")
	}

	if appointment.Status == appointments_models.CanceledByPatient {
		return fmt.Errorf("appointment status cannot change from %s to EDITED_BY_PSYCHOLOGIST", string(appointment.Status))
	}

	if input.Start < time.Now().Unix() {
		return errors.New("appointment cannot be scheduled to the past")
	}

	if input.Start >= input.End {
		return errors.New("appointment cannot have negative duration")
	}

	appointment.Status = appointments_models.EditedByPsychologist
	appointment.Start = input.Start
	appointment.End = input.End
	appointment.PriceRangeName = input.PriceRangeName
	appointment.Reason = input.Reason

	result = s.OrmUtil.Db().Save(&appointment)
	if result.Error != nil {
		return result.Error
	}

	return nil

}
