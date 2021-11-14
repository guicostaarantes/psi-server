package users_services

import (
	"errors"
	"time"

	users_models "github.com/guicostaarantes/psi-server/modules/users/models"
	"github.com/guicostaarantes/psi-server/utils/hash"
	"github.com/guicostaarantes/psi-server/utils/orm"
	"github.com/guicostaarantes/psi-server/utils/serializing"
	"github.com/guicostaarantes/psi-server/utils/token"
)

// AuthenticateUserService is a service that exchanges credentials for an access token
type AuthenticateUserService struct {
	HashUtil        hash.IHashUtil
	OrmUtil         orm.IOrmUtil
	SerializingUtil serializing.ISerializingUtil
	TokenUtil       token.ITokenUtil
	SecondsToExpire int64
}

// Execute is the method that runs the business logic of the service
func (s AuthenticateUserService) Execute(authInput *users_models.AuthenticateUserInput) (*users_models.Authentication, error) {

	user := users_models.User{}

	result := s.OrmUtil.Db().Where("email = ?", authInput.Email).Limit(1).Find(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	if user.ID == "" || !user.Active {
		return nil, errors.New("incorrect credentials")
	}

	compareErr := s.HashUtil.Compare(authInput.Password, user.Password)
	if compareErr != nil {
		if compareErr.Error() == s.HashUtil.GetWrongPasswordError() {
			return nil, errors.New("incorrect credentials")
		}
		return nil, compareErr
	}

	token, tokenErr := s.TokenUtil.GenerateToken(user.ID, s.SecondsToExpire)
	if tokenErr != nil {
		return nil, tokenErr
	}

	auth := &users_models.Authentication{
		UserID:    user.ID,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Second * time.Duration(s.SecondsToExpire)).Unix(),
		Token:     token,
	}

	result = s.OrmUtil.Db().Where("user_id = ?", auth.UserID).Delete(&users_models.Authentication{})
	if result.Error != nil {
		return nil, result.Error
	}

	result = s.OrmUtil.Db().Create(&auth)
	if result.Error != nil {
		return nil, result.Error
	}

	return auth, nil

}
