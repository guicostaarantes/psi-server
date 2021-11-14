package treatments_services

import (
	treatments_models "github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// GetTreatmentForPatientService is a service that gets a treatment based on its id
type GetTreatmentForPatientService struct {
	OrmUtil orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s GetTreatmentForPatientService) Execute(id string) (*treatments_models.GetPatientTreatmentsResponse, error) {

	treatment := &treatments_models.GetPatientTreatmentsResponse{}

	result := s.OrmUtil.Db().Model(&treatments_models.Treatment{}).Where("id = ?", id).Limit(1).Find(&treatment)
	if result.Error != nil {
		return nil, result.Error
	}

	return treatment, nil

}
