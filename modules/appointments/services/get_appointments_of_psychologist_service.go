package services

import (
	"github.com/guicostaarantes/psi-server/modules/appointments/models"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// GetAppointmentsOfPsychologistService is a service that the psychologist will use to retrieve their appointments
type GetAppointmentsOfPsychologistService struct {
	OrmUtil orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s GetAppointmentsOfPsychologistService) Execute(psychologistID string) ([]*models.Appointment, error) {

	appointments := []*models.Appointment{}

	result := s.OrmUtil.Db().Where("psychologist_id = ?", psychologistID).Order("created_at ASC").Find(&appointments)
	if result.Error != nil {
		return nil, result.Error
	}

	return appointments, nil

}
