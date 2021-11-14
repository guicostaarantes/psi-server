package appointments_services

import (
	appointments_models "github.com/guicostaarantes/psi-server/modules/appointments/models"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// GetAppointmentsOfPsychologistService is a service that the psychologist will use to retrieve their appointments
type GetAppointmentsOfPsychologistService struct {
	OrmUtil orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s GetAppointmentsOfPsychologistService) Execute(psychologistID string) ([]*appointments_models.Appointment, error) {

	appointments := []*appointments_models.Appointment{}

	result := s.OrmUtil.Db().Where("psychologist_id = ?", psychologistID).Order("created_at ASC").Find(&appointments)
	if result.Error != nil {
		return nil, result.Error
	}

	return appointments, nil

}
