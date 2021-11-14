package treatments_services

import (
	"errors"
	"fmt"
	"time"

	appointments_models "github.com/guicostaarantes/psi-server/modules/appointments/models"
	treatments_models "github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// InterruptTreatmentByPsychologistService is a service that interrupts a treatment, changing its status to interrupted by psychologist
type InterruptTreatmentByPsychologistService struct {
	OrmUtil orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s InterruptTreatmentByPsychologistService) Execute(id string, psychologistID string, reason string) error {

	treatment := treatments_models.Treatment{}

	result := s.OrmUtil.Db().Where("id = ? AND psychologist_id = ?", id, psychologistID).Limit(1).Find(&treatment)
	if result.Error != nil {
		return result.Error
	}

	if treatment.ID == "" {
		return errors.New("resource not found")
	}

	if treatment.Status != treatments_models.Active {
		return fmt.Errorf("treatments can only be interrupted if their current status is ACTIVE. current status is %s", string(treatment.Status))
	}

	appointments := []*appointments_models.Appointment{}

	result = s.OrmUtil.Db().Where("treatment_id = ? AND start > ?", id, time.Now().Unix()).Find(&appointments)
	if result.Error != nil {
		return result.Error
	}

	for _, appointment := range appointments {
		if appointment.Status != appointments_models.CanceledByPatient {
			appointment.Status = appointments_models.TreatmentInterruptedByPsychologist
			appointment.Reason = reason

			result = s.OrmUtil.Db().Save(&appointment)
			if result.Error != nil {
				return result.Error
			}
		}
	}

	treatment.EndDate = time.Now().Unix()
	treatment.Status = treatments_models.InterruptedByPsychologist
	treatment.Reason = reason

	result = s.OrmUtil.Db().Save(&treatment)
	if result.Error != nil {
		return result.Error
	}

	return nil

}
