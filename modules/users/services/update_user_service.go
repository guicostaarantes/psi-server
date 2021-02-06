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
func (s UpdateUserService) Execute(userID string, userInput *models.UpdateUserInput) error {

	user := models.User{}

	findErr := s.DatabaseUtil.FindOne("psi_db", "users", map[string]interface{}{"id": userID}, &user)
	if findErr != nil {
		return findErr
	}

	if user.ID == "" {
		return errors.New("resource not found")
	}

	user.FirstName = userInput.FirstName
	user.LastName = userInput.LastName
	user.Role = userInput.Role

	updateErr := s.DatabaseUtil.UpdateOne("psi_db", "users", map[string]interface{}{"id": userID}, user)
	if updateErr != nil {
		return updateErr
	}

	return nil

}
