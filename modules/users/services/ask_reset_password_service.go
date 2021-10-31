package services

import (
	"bytes"
	"html/template"
	"os"
	"time"

	mails_models "github.com/guicostaarantes/psi-server/modules/mails/models"
	models "github.com/guicostaarantes/psi-server/modules/users/models"
	"github.com/guicostaarantes/psi-server/modules/users/templates"
	"github.com/guicostaarantes/psi-server/utils/identifier"
	"github.com/guicostaarantes/psi-server/utils/orm"
	"github.com/guicostaarantes/psi-server/utils/token"
	"gorm.io/gorm"
)

// AskResetPasswordService is a service that sends a token to the user's email so that they can reset their password
type AskResetPasswordService struct {
	IdentifierUtil    identifier.IIdentifierUtil
	OrmUtil           orm.IOrmUtil
	TokenUtil         token.ITokenUtil
	SecondsToCooldown int64
	SecondsToExpire   int64
}

// Execute is the method that runs the business logic of the service
func (s AskResetPasswordService) Execute(email string) error {

	user := &models.User{}

	result := s.OrmUtil.Db().Where("email = ?", email).Limit(1).Find(&user)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return result.Error
	}

	if user.ID == "" || !user.Active {
		return nil
	}

	existsReset := &models.ResetPassword{}

	result = s.OrmUtil.Db().Where("user_id = ?", user.ID).Limit(1).Find(&existsReset)
	if result.Error != nil {
		return result.Error
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
		To:          user.Email,
		Cc:          "",
		Cco:         "",
		Subject:     "Redfinir senha do PSI",
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
