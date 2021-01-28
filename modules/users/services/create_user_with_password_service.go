package users_services

import (
	"errors"

	models "github.com/guicostaarantes/psi-server/modules/users/models"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/hash"
	"github.com/guicostaarantes/psi-server/utils/identifier"
	"github.com/guicostaarantes/psi-server/utils/match"
	"github.com/guicostaarantes/psi-server/utils/serializing"
)

// CreateUserWithPasswordService is a service that creates users and assigns them a password immediately
type CreateUserWithPasswordService struct {
	DatabaseUtil    database.IDatabaseUtil
	HashUtil        hash.IHashUtil
	IdentifierUtil  identifier.IIdentifierUtil
	MatchUtil       match.IMatchUtil
	SerializingUtil serializing.ISerializingUtil
}

// Execute is the method that runs the business logic of the service
func (s CreateUserWithPasswordService) Execute(userInput *models.CreateUserWithPasswordInput) error {

	emailErr := s.MatchUtil.IsEmailValid(userInput.Email)
	if emailErr != nil {
		return emailErr
	}

	passwordErr := s.MatchUtil.IsPasswordStrong(userInput.Password)
	if passwordErr != nil {
		return passwordErr
	}

	userWithSameEmail := models.User{}

	findErr := s.DatabaseUtil.FindOne("psi_db", "users", "email", userInput.Email, &userWithSameEmail)
	if findErr != nil {
		return findErr
	}

	if userWithSameEmail.ID != "" {
		return errors.New("user with same email already exists")
	}

	_, id, idErr := s.IdentifierUtil.GenerateIdentifier()
	if idErr != nil {
		return idErr
	}

	hashedPwd, hashErr := s.HashUtil.Hash(userInput.Password)
	if hashErr != nil {
		return hashErr
	}

	user := &models.User{
		ID:        string(id),
		Email:     userInput.Email,
		Active:    true,
		Password:  hashedPwd,
		FirstName: userInput.FirstName,
		LastName:  userInput.LastName,
		Role:      userInput.Role,
	}

	writeErr := s.DatabaseUtil.InsertOne("psi_db", "users", user)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
