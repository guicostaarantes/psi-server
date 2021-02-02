package users_services

import (
	"errors"

	models "github.com/guicostaarantes/psi-server/modules/users/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// GetUserByIdService is a service that gets the user by userId
type GetUserByIdService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s GetUserByIdService) Execute(id string) (*models.User, error) {

	user := &models.User{}

	findErr := s.DatabaseUtil.FindOne("psi_db", "users", map[string]interface{}{"id": id}, user)
	if findErr != nil {
		return nil, findErr
	}

	if user.ID == "" {
		return nil, errors.New("resource not found")
	}

	return user, nil

}
