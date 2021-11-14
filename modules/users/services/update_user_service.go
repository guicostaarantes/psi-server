package users_services

import (
	"errors"

	users_models "github.com/guicostaarantes/psi-server/modules/users/models"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// UpdateUserService is a service that can change data from a user
type UpdateUserService struct {
	OrmUtil orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s UpdateUserService) Execute(userID string, input *users_models.UpdateUserInput) error {

	user := users_models.User{}

	result := s.OrmUtil.Db().Where("id = ?", userID).Limit(1).Find(&user)
	if result.Error != nil {
		return result.Error
	}

	if user.ID == "" {
		return errors.New("resource not found")
	}

	user.Active = input.Active
	user.Role = input.Role

	result = s.OrmUtil.Db().Save(&user)
	if result.Error != nil {
		return result.Error
	}

	return nil

}
