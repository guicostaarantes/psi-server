package services

import (
	"github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// GetPsychologistPendingTreatmentsService is a service that gets all the pending treatments of a psychologist
type GetPsychologistPendingTreatmentsService struct {
	OrmUtil orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s GetPsychologistPendingTreatmentsService) Execute(psychologistID string) ([]*models.GetPsychologistTreatmentsResponse, error) {

	treatments := []*models.GetPsychologistTreatmentsResponse{}

	result := s.OrmUtil.Db().Model(&models.Treatment{}).Where("psychologist_id = ? AND status = ?", psychologistID, models.Pending).Order("created_at ASC").Find(&treatments)
	if result.Error != nil {
		return nil, result.Error
	}

	return treatments, nil

}
