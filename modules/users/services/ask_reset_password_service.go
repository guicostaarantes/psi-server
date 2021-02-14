package services

import (
	"bytes"
	"html/template"
	"os"
	"time"

	mails_models "github.com/guicostaarantes/psi-server/modules/mails/models"
	models "github.com/guicostaarantes/psi-server/modules/users/models"
	"github.com/guicostaarantes/psi-server/modules/users/templates"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/identifier"
	"github.com/guicostaarantes/psi-server/utils/token"
)

// AskResetPasswordService is a service that sends a token to the user's email so that they can reset their password
type AskResetPasswordService struct {
	DatabaseUtil      database.IDatabaseUtil
	IdentifierUtil    identifier.IIdentifierUtil
	TokenUtil         token.ITokenUtil
	SecondsToCooldown int64
	SecondsToExpire   int64
}

// Execute is the method that runs the business logic of the service
func (s AskResetPasswordService) Execute(email string) error {

	user := &models.User{}

	findUserErr := s.DatabaseUtil.FindOne("psi_db", "users", map[string]interface{}{"email": email}, user)
	if findUserErr != nil {
		return findUserErr
	}

	if user.ID == "" || !user.Active {
		return nil
	}

	existsReset := &models.ResetPassword{}

	findTokenErr := s.DatabaseUtil.FindOne("psi_db", "resets", map[string]interface{}{"userId": user.ID}, existsReset)
	if findTokenErr != nil {
		return findTokenErr
	}

	if existsReset.UserID != "" && existsReset.IssuedAt > time.Now().Add(-time.Second*time.Duration(s.SecondsToCooldown)).Unix() {
		return nil
	}

	token, tokenErr := s.TokenUtil.GenerateToken(user.ID, s.SecondsToCooldown)
	if tokenErr != nil {
		return tokenErr
	}

	reset := &models.ResetPassword{
		UserID:    user.ID,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Second * time.Duration(s.SecondsToCooldown)).Unix(),
		Token:     token,
		Redeemed:  false,
	}

	_, mailID, mailIDErr := s.IdentifierUtil.GenerateIdentifier()
	if mailIDErr != nil {
		return mailIDErr
	}

	templ, templErr := template.New("ResetPasswordEmail").Parse(templates.ResetPasswordEmailTemplate)
	if templErr != nil {
		return templErr
	}

	buff := new(bytes.Buffer)

	templ.Execute(buff, map[string]string{
		"SiteURL": os.Getenv("PSI_SITE_URL"),
		"Token":   token,
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

	return nil

}
