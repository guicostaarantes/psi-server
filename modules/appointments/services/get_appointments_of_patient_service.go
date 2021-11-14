package appointments_services

import (
	appointments_models "github.com/guicostaarantes/psi-server/modules/appointments/models"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// GetAppointmentsOfPatientService is a service that the patient will use to retrieve their appointments
type GetAppointmentsOfPatientService struct {
	OrmUtil orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s GetAppointmentsOfPatientService) Execute(patientID string) ([]*appointments_models.Appointment, error) {

	appointments := []*appointments_models.Appointment{}

	result := s.OrmUtil.Db().Where("patient_id = ?", patientID).Order("created_at ASC").Find(&appointments)
	if result.Error != nil {
		return nil, result.Error
	}

	return appointments, nil

}
