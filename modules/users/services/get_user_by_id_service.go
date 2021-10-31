package services

import (
	models "github.com/guicostaarantes/psi-server/modules/users/models"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// GetUserByIDService is a service that gets the user by userId
type GetUserByIDService struct {
	OrmUtil orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s GetUserByIDService) Execute(id string) (*models.User, error) {

	user := &models.User{}

	result := s.OrmUtil.Db().Where("id = ?", id).Limit(1).Find(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	if user.ID == "" {
		return nil, nil
	}

	return user, nil

}
