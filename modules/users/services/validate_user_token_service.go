package services

import (
	"errors"
	"time"

	models "github.com/guicostaarantes/psi-server/modules/users/models"
	"github.com/guicostaarantes/psi-server/utils/orm"
	"github.com/guicostaarantes/psi-server/utils/serializing"
)

// ValidateUserTokenService is a service that checks the validity of a token
type ValidateUserTokenService struct {
	OrmUtil         orm.IOrmUtil
	SerializingUtil serializing.ISerializingUtil
	SecondsToExpire int64
}

// Execute is the method that runs the business logic of the service
func (s ValidateUserTokenService) Execute(token string) (string, error) {

	auth := models.Authentication{}

	result := s.OrmUtil.Db().Where("token = ?", token).Limit(1).Find(&auth)
	if result.Error != nil {
		return "", result.Error
	}

	if auth.UserID == "" || auth.ExpiresAt < time.Now().Unix() {
		return "", errors.New("invalid token")
	}

	auth.ExpiresAt = time.Now().Unix() + s.SecondsToExpire

	result = s.OrmUtil.Db().Save(&auth)
	if result.Error != nil {
		return "", result.Error
	}

	return auth.UserID, nil

}
