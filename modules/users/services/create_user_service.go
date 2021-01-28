package users_services

import (
	"errors"

	mails_models "github.com/guicostaarantes/psi-server/modules/mails/models"
	models "github.com/guicostaarantes/psi-server/modules/users/models"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/identifier"
	"github.com/guicostaarantes/psi-server/utils/match"
	"github.com/guicostaarantes/psi-server/utils/serializing"
)

// CreateUserService is a service that creates users and sends emails so that the owner can assign a password
type CreateUserService struct {
	DatabaseUtil    database.IDatabaseUtil
	IdentifierUtil  identifier.IIdentifierUtil
	MatchUtil       match.IMatchUtil
	SerializingUtil serializing.ISerializingUtil
}

// Execute is the method that runs the business logic of the service
func (s CreateUserService) Execute(userInput *models.CreateUserInput) error {

	emailErr := s.MatchUtil.IsEmailValid(userInput.Email)
	if emailErr != nil {
		return emailErr
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

	user := &models.User{
		ID:        id,
		Email:     userInput.Email,
		Active:    true,
		FirstName: userInput.FirstName,
		LastName:  userInput.LastName,
		Role:      userInput.Role,
	}

	_, mailId, mailIdErr := s.IdentifierUtil.GenerateIdentifier()
	if mailIdErr != nil {
		return mailIdErr
	}

	mail := &mails_models.TransientMailMessage{
		ID:          mailId,
		FromAddress: "relacionamento@psi.com.br",
		FromName:    "Relacionamento PSI",
		To:          []string{userInput.Email},
		Cc:          []string{},
		Cco:         []string{},
		Subject:     "Hello",
		Html:        "Hello there",
		Processed:   false,
	}

	writeMailErr := s.DatabaseUtil.InsertOne("psi_db", "mails", mail)
	if writeMailErr != nil {
		return writeMailErr
	}

	writeErr := s.DatabaseUtil.InsertOne("psi_db", "users", user)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
