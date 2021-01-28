package users_services

import (
	models "github.com/guicostaarantes/psi-server/modules/users/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// ActivateUserService is a service that can activate or deactivate a user
type ActivateUserService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (v ActivateUserService) Execute(userID string, active bool) error {

	user := models.User{}

	findErr := v.DatabaseUtil.FindOne("psi_db", "users", "id", userID, &user)
	if findErr != nil {
		return findErr
	}

	user.Active = active

	updateErr := v.DatabaseUtil.UpdateOne("psi_db", "auths", "id", userID, user)
	if updateErr != nil {
		return updateErr
	}

	return nil

}
