package services

import (
	"bytes"
	"errors"
	"html/template"
	"os"
	"time"

	mails_models "github.com/guicostaarantes/psi-server/modules/mails/models"
	models "github.com/guicostaarantes/psi-server/modules/users/models"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/identifier"
	"github.com/guicostaarantes/psi-server/utils/match"
	"github.com/guicostaarantes/psi-server/utils/merge"
	"github.com/guicostaarantes/psi-server/utils/serializing"
	"github.com/guicostaarantes/psi-server/utils/token"
)

// CreateUserService is a service that creates users and sends emails so that the owner can assign a password
type CreateUserService struct {
	DatabaseUtil    database.IDatabaseUtil
	IdentifierUtil  identifier.IIdentifierUtil
	MatchUtil       match.IMatchUtil
	MergeUtil       merge.IMergeUtil
	SerializingUtil serializing.ISerializingUtil
	TokenUtil       token.ITokenUtil
	SecondsToExpire int64
}

// Execute is the method that runs the business logic of the service
func (s CreateUserService) Execute(userInput *models.CreateUserInput) error {

	emailErr := s.MatchUtil.IsEmailValid(userInput.Email)
	if emailErr != nil {
		return emailErr
	}

	userWithSameEmail := models.User{}

	findErr := s.DatabaseUtil.FindOne("psi_db", "users", map[string]interface{}{"email": userInput.Email}, &userWithSameEmail)
	if findErr != nil {
		return findErr
	}

	if userWithSameEmail.ID != "" {
		return errors.New("user with same email already exists")
	}

	_, userID, userIDErr := s.IdentifierUtil.GenerateIdentifier()
	if userIDErr != nil {
		return userIDErr
	}

	user := &models.User{
		ID:     userID,
		Active: true,
	}

	mergeErr := s.MergeUtil.Merge(&user, userInput)
	if mergeErr != nil {
		return mergeErr
	}

	token, tokenErr := s.TokenUtil.GenerateToken(user.ID, s.SecondsToExpire)
	if tokenErr != nil {
		return tokenErr
	}

	reset := &models.ResetPassword{
		UserID:    userID,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Second * time.Duration(s.SecondsToExpire)).Unix(),
		Token:     token,
		Redeemed:  false,
	}

	_, mailID, mailIDErr := s.IdentifierUtil.GenerateIdentifier()
	if mailIDErr != nil {
		return mailIDErr
	}

	templ, templErr := template.ParseFiles("templates/create_user_email.html")
	if templErr != nil {
		return templErr
	}

	buff := new(bytes.Buffer)

	templ.Execute(buff, map[string]string{
		"FirstName": user.FirstName,
		"SiteURL":   os.Getenv("PSI_SITE_URL"),
		"Token":     token,
	})

	mail := &mails_models.TransientMailMessage{
		ID:          mailID,
		FromAddress: "relacionamento@psi.com.br",
		FromName:    "Relacionamento PSI",
		To:          []string{user.Email},
		Cc:          []string{},
		Cco:         []string{},
		Subject:     "Bem-vindo ao PSI",
		Html:        buff.String(),
		Processed:   false,
	}

	writeMailErr := s.DatabaseUtil.InsertOne("psi_db", "mails", mail)
	if writeMailErr != nil {
		return writeMailErr
	}

	writeResetErr := s.DatabaseUtil.InsertOne("psi_db", "resets", reset)
	if writeResetErr != nil {
		return writeResetErr
	}

	writeUserErr := s.DatabaseUtil.InsertOne("psi_db", "users", user)
	if writeUserErr != nil {
		return writeUserErr
	}

	return nil

}
