package services

import (
	models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// GetPsychologistService is a service that gets the psychologist profile based on id
type GetPsychologistService struct {
	OrmUtil orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s GetPsychologistService) Execute(id string) (*models.Psychologist, error) {

	psy := &models.Psychologist{}

	result := s.OrmUtil.Db().Where("id = ?", id).Limit(1).Find(&psy)
	if result.Error != nil {
		return nil, result.Error
	}

	if psy.ID == "" {
		return nil, nil
	}

	return psy, nil

}
