package services

import (
	"github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// GetTreatmentForPsychologistService is a service that gets a treatment based on its id
type GetTreatmentForPsychologistService struct {
	OrmUtil orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s GetTreatmentForPsychologistService) Execute(id string) (*models.GetPsychologistTreatmentsResponse, error) {

	treatment := &models.GetPsychologistTreatmentsResponse{}

	result := s.OrmUtil.Db().Model(&models.Treatment{}).Where("id = ?", id).Limit(1).Find(&treatment)
	if result.Error != nil {
		return nil, result.Error
	}

	return treatment, nil

}
