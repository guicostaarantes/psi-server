package users_services

import (
	"errors"
	"time"

	users_models "github.com/guicostaarantes/psi-server/modules/users/models"
	"github.com/guicostaarantes/psi-server/utils/orm"
	"github.com/guicostaarantes/psi-server/utils/serializing"
)

// ValidateUserTokenService is a service that checks the validity of a token
type ValidateUserTokenService struct {
	OrmUtil                 orm.IOrmUtil
	SerializingUtil         serializing.ISerializingUtil
	ExpireAuthTokenDuration time.Duration
}

// Execute is the method that runs the business logic of the service
func (s ValidateUserTokenService) Execute(token string) (string, error) {

	auth := users_models.Authentication{}

	result := s.OrmUtil.Db().Where("token = ?", token).Limit(1).Find(&auth)
	if result.Error != nil {
		return "", result.Error
	}

	if auth.UserID == "" || auth.ExpiresAt.Before(time.Now()) {
		return "", errors.New("invalid token")
	}

	auth.ExpiresAt = time.Now().Add(s.ExpireAuthTokenDuration)

	result = s.OrmUtil.Db().Save(&auth)
	if result.Error != nil {
		return "", result.Error
	}

	return auth.UserID, nil

}
