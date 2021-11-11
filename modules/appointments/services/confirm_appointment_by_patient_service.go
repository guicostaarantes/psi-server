package services

import (
	"errors"
	"fmt"

	"github.com/guicostaarantes/psi-server/modules/appointments/models"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// ConfirmAppointmentByPatientService is a service that the patient will use to confirm an appointment
type ConfirmAppointmentByPatientService struct {
	OrmUtil orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s ConfirmAppointmentByPatientService) Execute(id string, patientID string) error {

	appointment := models.Appointment{}

	result := s.OrmUtil.Db().Where("id = ? AND patient_id = ?", id, patientID).Limit(1).Find(&appointment)
	if result.Error != nil {
		return result.Error
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

	result = s.OrmUtil.Db().Save(&appointment)
	if result.Error != nil {
		return result.Error
	}

	return nil

}
