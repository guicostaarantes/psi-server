package services

import (
	"errors"
	"time"

	models "github.com/guicostaarantes/psi-server/modules/users/models"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/serializing"
)

// ValidateUserTokenService is a service that checks the validity of a token
type ValidateUserTokenService struct {
	DatabaseUtil    database.IDatabaseUtil
	SerializingUtil serializing.ISerializingUtil
	SecondsToExpire int64
}

// Execute is the method that runs the business logic of the service
func (s ValidateUserTokenService) Execute(token string) (string, error) {

	auth := models.Authentication{}

	findErr := s.DatabaseUtil.FindOne("auths", map[string]interface{}{"token": token}, &auth)
	if findErr != nil {
		return "", findErr
	}

	if auth.UserID == "" || auth.ExpiresAt < time.Now().Unix() {
		return "", errors.New("invalid token")
	}

	auth.ExpiresAt = time.Now().Unix() + s.SecondsToExpire

	updateErr := s.DatabaseUtil.UpdateOne("auths", map[string]interface{}{"token": token}, auth)
	if updateErr != nil {
		return "", updateErr
	}

	return auth.UserID, nil

}
