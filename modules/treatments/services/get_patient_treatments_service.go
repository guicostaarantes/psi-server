package services

import (
	"github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// GetPatientTreatmentsService is a service that gets all the treatments of a psychologist
type GetPatientTreatmentsService struct {
	OrmUtil orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s GetPatientTreatmentsService) Execute(patientID string) ([]*models.GetPatientTreatmentsResponse, error) {

	treatments := []*models.GetPatientTreatmentsResponse{}

	result := s.OrmUtil.Db().Model(&models.Treatment{}).Where("patient_id = ?", patientID).Order("created_at ASC").Find(&treatments)
	if result.Error != nil {
		return nil, result.Error
	}

	return treatments, nil

}
