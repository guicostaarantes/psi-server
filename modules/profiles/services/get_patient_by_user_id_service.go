package services

import (
	models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// GetPatientByUserIDService is a service that gets the patient profile based on UserID
type GetPatientByUserIDService struct {
	OrmUtil orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s GetPatientByUserIDService) Execute(id string) (*models.Patient, error) {

	patient := &models.Patient{}

	result := s.OrmUtil.Db().Where("user_id = ?", id).Limit(1).Find(&patient)
	if result.Error != nil {
		return nil, result.Error
	}

	if patient.ID == "" {
		return nil, nil
	}

	return patient, nil

}
