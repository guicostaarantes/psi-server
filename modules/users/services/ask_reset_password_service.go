package users_services

import (
	"bytes"
	"html/template"
	"os"
	"time"

	mails_models "github.com/guicostaarantes/psi-server/modules/mails/models"
	users_models "github.com/guicostaarantes/psi-server/modules/users/models"
	users_templates "github.com/guicostaarantes/psi-server/modules/users/templates"
	"github.com/guicostaarantes/psi-server/utils/identifier"
	"github.com/guicostaarantes/psi-server/utils/orm"
	"github.com/guicostaarantes/psi-server/utils/token"
)

// AskResetPasswordService is a service that sends a token to the user's email so that they can reset their password
type AskResetPasswordService struct {
	IdentifierUtil           identifier.IIdentifierUtil
	OrmUtil                  orm.IOrmUtil
	TokenUtil                token.ITokenUtil
	ExpireResetTokenDuration time.Duration
}

// Execute is the method that runs the business logic of the service
func (s AskResetPasswordService) Execute(email string) error {

	user := &users_models.User{}

	result := s.OrmUtil.Db().Where("email = ?", email).Limit(1).Find(&user)
	if result.Error != nil {
		return result.Error
	}

	if user.ID == "" || !user.Active {
		return nil
	}

	existingReset := &users_models.ResetPassword{}

	result = s.OrmUtil.Db().Where("user_id = ?", user.ID).Limit(1).Find(&existingReset)
	if result.Error != nil {
		return result.Error
	}

	if existingReset.UserID != "" {
		if existingReset.ExpiresAt.After(time.Now()) {
			return nil
		}

		result := s.OrmUtil.Db().Delete(&existingReset)
		if result.Error != nil {
			return result.Error
		}
	}

	token, tokenErr := s.TokenUtil.GenerateToken(user.ID, s.ExpireResetTokenDuration)
	if tokenErr != nil {
		return tokenErr
	}

	reset := &users_models.ResetPassword{
		UserID:    user.ID,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(s.ExpireResetTokenDuration),
		Token:     token,
		Redeemed:  false,
	}

	_, mailID, mailIDErr := s.IdentifierUtil.GenerateIdentifier()
	if mailIDErr != nil {
		return mailIDErr
	}

	templ, templErr := template.New("ResetPasswordEmail").Parse(users_templates.ResetPasswordEmailTemplate)
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
		To:          user.Email,
		Cc:          "",
		Cco:         "",
		Subject:     "Redefinir senha do PSI",
		Html:        buff.String(),
		Processed:   false,
	}

	result = s.OrmUtil.Db().Create(&mail)
	if result.Error != nil {
		return result.Error
	}

	result = s.OrmUtil.Db().Create(&reset)
	if result.Error != nil {
		return result.Error
	}

	return nil

}
