package services

import (
	"errors"

	models "github.com/guicostaarantes/psi-server/modules/users/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// UpdateUserService is a service that can change data from a user
type UpdateUserService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s UpdateUserService) Execute(userID string, input *models.UpdateUserInput) error {

	user := models.User{}

	findErr := s.DatabaseUtil.FindOne("users", map[string]interface{}{"id": userID}, &user)
	if findErr != nil {
		return findErr
	}

	if user.ID == "" {
		return errors.New("resource not found")
	}

	user.Active = input.Active
	user.Role = input.Role

	updateErr := s.DatabaseUtil.UpdateOne("users", map[string]interface{}{"id": userID}, user)
	if updateErr != nil {
		return updateErr
	}

	return nil

}
