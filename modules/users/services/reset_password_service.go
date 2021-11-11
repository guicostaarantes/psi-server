package services

import (
	"errors"
	"time"

	models "github.com/guicostaarantes/psi-server/modules/users/models"
	"github.com/guicostaarantes/psi-server/utils/hash"
	"github.com/guicostaarantes/psi-server/utils/match"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// ResetPasswordService is a service that (re)sets the password for a user based on a token sent to their email
type ResetPasswordService struct {
	MatchUtil match.IMatchUtil
	HashUtil  hash.IHashUtil
	OrmUtil   orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s ResetPasswordService) Execute(resetInput *models.ResetPasswordInput) error {

	passwordErr := s.MatchUtil.IsPasswordStrong(resetInput.Password)
	if passwordErr != nil {
		return passwordErr
	}

	reset := &models.ResetPassword{}

	result := s.OrmUtil.Db().Where("token = ?", resetInput.Token).Limit(1).Find(&reset)
	if result.Error != nil {
		return result.Error
	}

	if reset.UserID == "" || reset.ExpiresAt < time.Now().Unix() || reset.Redeemed {
		return errors.New("invalid token")
	}

	user := &models.User{}

	result = s.OrmUtil.Db().Where("id = ?", reset.UserID).Limit(1).Find(&user)
	if result.Error != nil {
		return result.Error
	}

	if user.ID == "" || !user.Active {
		return errors.New("invalid token")
	}

	hashedPwd, hashErr := s.HashUtil.Hash(resetInput.Password)
	if hashErr != nil {
		return hashErr
	}

	user.Password = hashedPwd

	result = s.OrmUtil.Db().Save(&user)
	if result.Error != nil {
		return result.Error
	}

	reset.Redeemed = true

	result = s.OrmUtil.Db().Save(&reset)
	if result.Error != nil {
		return result.Error
	}

	return nil

}
