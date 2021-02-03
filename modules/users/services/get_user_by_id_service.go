package services

import (
	"errors"

	models "github.com/guicostaarantes/psi-server/modules/users/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// GetUserByIDService is a service that gets the user by userId
type GetUserByIDService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s GetUserByIDService) Execute(id string) (*models.User, error) {

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
