package services

import (
	"errors"

	models "github.com/guicostaarantes/psi-server/modules/users/models"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/merge"
)

// UpdateUserService is a service that can change data from a user
type UpdateUserService struct {
	DatabaseUtil database.IDatabaseUtil
	MergeUtil    merge.IMergeUtil
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

	mergeErr := s.MergeUtil.Merge(&user, userInput)
	if mergeErr != nil {
		return mergeErr
	}

	updateErr := s.DatabaseUtil.UpdateOne("psi_db", "users", map[string]interface{}{"id": userID}, user)
	if updateErr != nil {
		return updateErr
	}

	return nil

}
