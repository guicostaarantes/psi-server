package treatments_services

import (
	treatments_models "github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// GetPsychologistTreatmentsService is a service that gets all the treatments of a psychologist
type GetPsychologistTreatmentsService struct {
	OrmUtil orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s GetPsychologistTreatmentsService) Execute(psychologistID string) ([]*treatments_models.GetPsychologistTreatmentsResponse, error) {

	treatments := []*treatments_models.GetPsychologistTreatmentsResponse{}

	result := s.OrmUtil.Db().Model(&treatments_models.Treatment{}).Where("psychologist_id = ?", psychologistID).Order("created_at ASC").Find(&treatments)
	if result.Error != nil {
		return nil, result.Error
	}

	return treatments, nil

}
