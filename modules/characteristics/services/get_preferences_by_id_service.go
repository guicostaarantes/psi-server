package services

import (
	"github.com/guicostaarantes/psi-server/modules/characteristics/models"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// GetPreferencesByIDService is a service that gets the preferences of a profile based on its id
type GetPreferencesByIDService struct {
	OrmUtil orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s GetPreferencesByIDService) Execute(id string) ([]*models.PreferenceResponse, error) {

	preferences := []*models.PreferenceResponse{}

	result := s.OrmUtil.Db().Model(&models.Preference{}).Where("profile_id = ?", id).Find(&preferences)
	if result.Error != nil {
		return nil, result.Error
	}

	return preferences, nil

}
