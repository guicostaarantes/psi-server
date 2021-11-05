package services

import (
	"errors"
	"fmt"
	"time"

	appointments_models "github.com/guicostaarantes/psi-server/modules/appointments/models"
	"github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// FinalizeTreatmentService is a service that changes the status of a treatment to finalized
type FinalizeTreatmentService struct {
	OrmUtil orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s FinalizeTreatmentService) Execute(id string, psychologistID string) error {

	treatment := models.Treatment{}

	result := s.OrmUtil.Db().Where("id = ? AND psychologist_id = ?", id, psychologistID).Limit(1).Find(&treatment)
	if result.Error != nil {
		return result.Error
	}

	if treatment.ID == "" {
		return errors.New("resource not found")
	}

	if treatment.Status != models.Active {
		return fmt.Errorf("treatments can only be finalized if their current status is ACTIVE. current status is %s", string(treatment.Status))
	}

	appointmentsOfTreatment := []*appointments_models.Appointment{}

	result = s.OrmUtil.Db().Where("treatment_id = ?", treatment.ID).Find(&appointmentsOfTreatment)
	if result.Error != nil {
		return result.Error
	}

	for _, appointment := range appointmentsOfTreatment {
		if appointment.Start > time.Now().Unix() && appointment.Status != appointments_models.CanceledByPatient {
			appointment.Status = appointments_models.TreatmentFinalized
			appointment.Reason = "Tratamento finalizado"

			result = s.OrmUtil.Db().Save(&appointment)
			if result.Error != nil {
				return result.Error
			}
		}
	}

	treatment.EndDate = time.Now().Unix()
	treatment.Status = models.Finalized

	result = s.OrmUtil.Db().Save(&treatment)
	if result.Error != nil {
		return result.Error
	}

	return nil

}
