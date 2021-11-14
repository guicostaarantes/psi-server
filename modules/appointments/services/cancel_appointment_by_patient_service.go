package appointments_services

import (
	"errors"
	"fmt"

	appointments_models "github.com/guicostaarantes/psi-server/modules/appointments/models"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// CancelAppointmentByPatientService is a service that the patient will use to cancel an appointment
type CancelAppointmentByPatientService struct {
	OrmUtil orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s CancelAppointmentByPatientService) Execute(id string, patientID string, reason string) error {

	appointment := appointments_models.Appointment{}

	result := s.OrmUtil.Db().Where("id = ? AND patient_id = ?", id, patientID).Limit(1).Find(&appointment)
	if result.Error != nil {
		return result.Error
	}

	if appointment.ID == "" {
		return errors.New("resource not found")
	}

	if appointment.Status == appointments_models.CanceledByPatient || appointment.Status == appointments_models.CanceledByPsychologist {
		return fmt.Errorf("appointment status cannot change from %s to CANCELED_BY_PATIENT", string(appointment.Status))
	}

	appointment.Status = appointments_models.CanceledByPatient
	appointment.Reason = reason

	result = s.OrmUtil.Db().Save(&appointment)
	if result.Error != nil {
		return result.Error
	}

	return nil

}
