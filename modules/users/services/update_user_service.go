package users_services

import (
	"fmt"

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
func (u UpdateUserService) Execute(userID string, userInput *models.UpdateUserInput) error {

	user := models.User{}

	findErr := u.DatabaseUtil.FindOne("psi_db", "users", "id", userID, &user)
	if findErr != nil {
		return findErr
	}

	mergeErr := u.MergeUtil.Merge(&user, userInput)
	if mergeErr != nil {
		return mergeErr
	}

	fmt.Printf("%#v \n", userInput)

	updateErr := u.DatabaseUtil.UpdateOne("psi_db", "users", "id", userID, user)
	if updateErr != nil {
		return updateErr
	}

	return nil

}
