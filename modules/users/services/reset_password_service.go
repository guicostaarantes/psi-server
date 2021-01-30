package users_services

import (
	"errors"
	"time"

	models "github.com/guicostaarantes/psi-server/modules/users/models"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/hash"
	"github.com/guicostaarantes/psi-server/utils/match"
)

// ResetPasswordService is a service that (re)sets the password for a user based on a token sent to their email
type ResetPasswordService struct {
	DatabaseUtil database.IDatabaseUtil
	MatchUtil    match.IMatchUtil
	HashUtil     hash.IHashUtil
}

// Execute is the method that runs the business logic of the service
func (r ResetPasswordService) Execute(resetInput *models.ResetPasswordInput) error {

	passwordErr := r.MatchUtil.IsPasswordStrong(resetInput.Password)
	if passwordErr != nil {
		return passwordErr
	}

	reset := &models.ResetPassword{}

	findTokenErr := r.DatabaseUtil.FindOne("psi_db", "resets", "token", resetInput.Token, reset)
	if findTokenErr != nil {
		return findTokenErr
	}

	if reset.UserID == "" || reset.ExpiresAt < time.Now().Unix() {
		return errors.New("invalid token")
	}

	user := &models.User{}

	findUserErr := r.DatabaseUtil.FindOne("psi_db", "users", "id", reset.UserID, user)
	if findUserErr != nil {
		return findUserErr
	}

	if user.ID == "" || !user.Active {
		return errors.New("invalid token")
	}

	hashedPwd, hashErr := r.HashUtil.Hash(resetInput.Password)
	if hashErr != nil {
		return hashErr
	}

	user.Password = hashedPwd

	updateUserErr := r.DatabaseUtil.UpdateOne("psi_db", "users", "id", reset.UserID, user)
	if updateUserErr != nil {
		return updateUserErr
	}

	deleteTokenErr := r.DatabaseUtil.DeleteOne("psi_db", "resets", "token", resetInput.Token)
	if deleteTokenErr != nil {
		return deleteTokenErr
	}

	return nil

}
