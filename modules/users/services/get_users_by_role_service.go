package services

import (
	models "github.com/guicostaarantes/psi-server/modules/users/models"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// GetUsersByRoleService is a service that gets the all users of a specific role in the database
type GetUsersByRoleService struct {
	OrmUtil orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s GetUsersByRoleService) Execute(role string) ([]*models.User, error) {

	users := []*models.User{}

	result := s.OrmUtil.Db().Where("role = ?", role).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil

}
