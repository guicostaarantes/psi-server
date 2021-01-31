package users_services

import (
	"errors"
	"time"

	models "github.com/guicostaarantes/psi-server/modules/users/models"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/hash"
	"github.com/guicostaarantes/psi-server/utils/serializing"
	"github.com/guicostaarantes/psi-server/utils/token"
)

// AuthenticateUserService is a service that exchanges credentials for an access token
type AuthenticateUserService struct {
	DatabaseUtil    database.IDatabaseUtil
	HashUtil        hash.IHashUtil
	SerializingUtil serializing.ISerializingUtil
	TokenUtil       token.ITokenUtil
	SecondsToExpire int64
}

// Execute is the method that runs the business logic of the service
func (s AuthenticateUserService) Execute(authInput *models.AuthenticateUserInput) (*models.Authentication, error) {

	if authInput.IPAddress == "" {
		return nil, errors.New("must notify IP")
	}

	user := models.User{}

	findErr := s.DatabaseUtil.FindOne("psi_db", "users", "email", authInput.Email, &user)
	if findErr != nil {
		return nil, findErr
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

	auth := &models.Authentication{
		UserID:    user.ID,
		IPAddress: authInput.IPAddress,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Second * time.Duration(s.SecondsToExpire)).Unix(),
		Token:     token,
	}

	deleteErr := s.DatabaseUtil.DeleteOne("psi_db", "auths", "userId", auth.UserID)
	if deleteErr != nil {
		return nil, deleteErr
	}

	writeErr := s.DatabaseUtil.InsertOne("psi_db", "auths", auth)
	if writeErr != nil {
		return nil, writeErr
	}

	return auth, nil

}
