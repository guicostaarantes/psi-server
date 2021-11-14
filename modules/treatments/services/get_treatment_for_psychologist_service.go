package treatments_services

import (
	treatments_models "github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// GetTreatmentForPsychologistService is a service that gets a treatment based on its id
type GetTreatmentForPsychologistService struct {
	OrmUtil orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s GetTreatmentForPsychologistService) Execute(id string) (*treatments_models.GetPsychologistTreatmentsResponse, error) {

	treatment := &treatments_models.GetPsychologistTreatmentsResponse{}

	result := s.OrmUtil.Db().Model(&treatments_models.Treatment{}).Where("id = ?", id).Limit(1).Find(&treatment)
	if result.Error != nil {
		return nil, result.Error
	}

	return treatment, nil

}
