package users_services

import (
	"bytes"
	"errors"
	"html/template"
	"os"
	"time"

	mails_models "github.com/guicostaarantes/psi-server/modules/mails/models"
	users_models "github.com/guicostaarantes/psi-server/modules/users/models"
	users_templates "github.com/guicostaarantes/psi-server/modules/users/templates"
	"github.com/guicostaarantes/psi-server/utils/identifier"
	"github.com/guicostaarantes/psi-server/utils/match"
	"github.com/guicostaarantes/psi-server/utils/orm"
	"github.com/guicostaarantes/psi-server/utils/serializing"
	"github.com/guicostaarantes/psi-server/utils/token"
)

// CreateUserService is a service that creates users and sends emails so that the owner can assign a password
type CreateUserService struct {
	IdentifierUtil          identifier.IIdentifierUtil
	MatchUtil               match.IMatchUtil
	OrmUtil                 orm.IOrmUtil
	SerializingUtil         serializing.ISerializingUtil
	TokenUtil               token.ITokenUtil
	ExpireAuthTokenDuration time.Duration
}

// Execute is the method that runs the business logic of the service
func (s CreateUserService) Execute(userInput *users_models.CreateUserInput, whichTemplate string) error {

	emailErr := s.MatchUtil.IsEmailValid(userInput.Email)
	if emailErr != nil {
		return emailErr
	}

	userWithSameEmail := users_models.User{}

	result := s.OrmUtil.Db().Where("email = ?", userInput.Email).Limit(1).Find(&userWithSameEmail)
	if result.Error != nil {
		return result.Error
	}

	// If user with same email exists, will not send email but will succeed in order not to inform hackers what emails are in the system
	if userWithSameEmail.ID != "" {
		return nil
	}

	_, userID, userIDErr := s.IdentifierUtil.GenerateIdentifier()
	if userIDErr != nil {
		return userIDErr
	}

	user := &users_models.User{
		ID:     userID,
		Active: true,
		Email:  userInput.Email,
		Role:   userInput.Role,
	}

	token, tokenErr := s.TokenUtil.GenerateToken(user.ID, s.ExpireAuthTokenDuration)
	if tokenErr != nil {
		return tokenErr
	}

	reset := &users_models.ResetPassword{
		UserID:    userID,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(s.ExpireAuthTokenDuration),
		Token:     token,
		Redeemed:  false,
	}

	_, mailID, mailIDErr := s.IdentifierUtil.GenerateIdentifier()
	if mailIDErr != nil {
		return mailIDErr
	}

	templateToUse := ""

	switch whichTemplate {
	case "CREATE_PATIENT":
		templateToUse = users_templates.CreateUserEmailTemplate
	case "INVITE_PSYCHOLOGIST":
		templateToUse = users_templates.InvitePsychologistEmailTemplate
	default:
		return errors.New("template not found")
	}

	templ, templErr := template.New("CreateUserEmail").Parse(templateToUse)
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
		Subject:     "Bem-vindo ao PSI",
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

	result = s.OrmUtil.Db().Create(&user)
	if result.Error != nil {
		return result.Error
	}

	return nil

}
