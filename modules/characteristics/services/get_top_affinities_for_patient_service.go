package characteristcs_services

import (
	"fmt"
	"time"

	characteristics_models "github.com/guicostaarantes/psi-server/modules/characteristics/models"
	cooldowns_models "github.com/guicostaarantes/psi-server/modules/cooldowns/models"
	cooldowns_services "github.com/guicostaarantes/psi-server/modules/cooldowns/services"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// GetTopAffinitiesForPatientService is a service that sets the top affinities for a specific patient profile if the cache is old enough, and returns them
type GetTopAffinitiesForPatientService struct {
	OrmUtil                           orm.IOrmUtil
	GetCooldownService                *cooldowns_services.GetCooldownService
	SetTopAffinitiesForPatientService *SetTopAffinitiesForPatientService
}

// Execute is the method that runs the business logic of the service
func (s GetTopAffinitiesForPatientService) Execute(patientID string) ([]*characteristics_models.Affinity, error) {
	cooldown, getErr := s.GetCooldownService.Execute(patientID, cooldowns_models.Patient, cooldowns_models.TreatmentInterrupted)
	if getErr != nil {
		return nil, getErr
	}

	if cooldown != nil {
		return nil, fmt.Errorf("assign treatment is blocked for this user until %s", cooldown.ValidUntil.Format(time.RFC3339))
	}

	cooldown, getErr = s.GetCooldownService.Execute(patientID, cooldowns_models.Patient, cooldowns_models.TopAffinitiesSet)
	if getErr != nil {
		return nil, getErr
	}

	if cooldown == nil {
		setErr := s.SetTopAffinitiesForPatientService.Execute(patientID)
		if setErr != nil {
			return nil, setErr
		}
	}

	topAffinities := []*characteristics_models.Affinity{}

	result := s.OrmUtil.Db().Where("patient_id = ?", patientID).Find(&topAffinities)
	if result.Error != nil {
		return nil, result.Error
	}

	return topAffinities, nil
}
