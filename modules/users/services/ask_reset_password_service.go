package users_services

import (
	"bytes"
	"html/template"
	"os"
	"time"

	mails_models "github.com/guicostaarantes/psi-server/modules/mails/models"
	models "github.com/guicostaarantes/psi-server/modules/users/models"
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
func (a AskResetPasswordService) Execute(email string) error {

	user := &models.User{}

	findUserErr := a.DatabaseUtil.FindOne("psi_db", "users", "email", email, user)
	if findUserErr != nil {
		return findUserErr
	}

	if user.ID == "" || !user.Active {
		return nil
	}

	existsReset := &models.ResetPassword{}

	findTokenErr := a.DatabaseUtil.FindOne("psi_db", "resets", "userId", user.ID, existsReset)
	if findTokenErr != nil {
		return findTokenErr
	}

	if existsReset.ID != "" && existsReset.IssuedAt > time.Now().Add(-time.Second*time.Duration(a.SecondsToCooldown)).Unix() {
		return nil
	}

	_, resetTokenID, resetTokenIDErr := a.IdentifierUtil.GenerateIdentifier()
	if resetTokenIDErr != nil {
		return resetTokenIDErr
	}

	token, tokenErr := a.TokenUtil.GenerateToken(user.ID, a.SecondsToCooldown)
	if tokenErr != nil {
		return tokenErr
	}

	reset := &models.ResetPassword{
		ID:        resetTokenID,
		UserID:    user.ID,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Second * time.Duration(a.SecondsToCooldown)).Unix(),
		Token:     token,
		Redeemed:  false,
	}

	_, mailID, mailIDErr := a.IdentifierUtil.GenerateIdentifier()
	if mailIDErr != nil {
		return mailIDErr
	}

	templ, templErr := template.ParseFiles("templates/reset_password_email.html")
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

	writeMailErr := a.DatabaseUtil.InsertOne("psi_db", "mails", mail)
	if writeMailErr != nil {
		return writeMailErr
	}

	writeResetErr := a.DatabaseUtil.InsertOne("psi_db", "resets", reset)
	if writeResetErr != nil {
		return writeResetErr
	}

	return nil

}
