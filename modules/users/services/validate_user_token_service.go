package users_services

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
func (v ValidateUserTokenService) Execute(token string) (string, error) {

	auth := models.Authentication{}

	findErr := v.DatabaseUtil.FindOne("psi_db", "auths", "token", token, &auth)
	if findErr != nil {
		return "", findErr
	}

	if auth.ID == "" || auth.ExpiresAt < time.Now().Unix() {
		return "", errors.New("invalid token")
	}

	auth.ExpiresAt = time.Now().Unix() + v.SecondsToExpire

	updateErr := v.DatabaseUtil.UpdateOne("psi_db", "auths", "token", token, auth)
	if updateErr != nil {
		return "", updateErr
	}

	return auth.UserID, nil

}
