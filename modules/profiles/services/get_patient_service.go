package services

import (
	models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// GetPatientService is a service that gets the patient profile based on id
type GetPatientService struct {
	OrmUtil orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s GetPatientService) Execute(id string) (*models.Patient, error) {

	patient := &models.Patient{}

	result := s.OrmUtil.Db().Where("id = ?", id).Limit(1).Find(&patient)
	if result.Error != nil {
		return nil, result.Error
	}

	if patient.ID == "" {
		return nil, nil
	}

	return patient, nil

}
