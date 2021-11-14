package services

import (
	profiles_models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// GetPsychologistByUserIDService is a service that gets the psychologist profile based on UserID
type GetPsychologistByUserIDService struct {
	OrmUtil orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s GetPsychologistByUserIDService) Execute(id string) (*profiles_models.Psychologist, error) {

	psy := &profiles_models.Psychologist{}

	result := s.OrmUtil.Db().Where("user_id = ?", id).Limit(1).Find(&psy)
	if result.Error != nil {
		return nil, result.Error
	}

	if psy.ID == "" {
		return nil, nil
	}

	return psy, nil

}
