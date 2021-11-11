package services

import (
	"errors"

	models "github.com/guicostaarantes/psi-server/modules/users/models"
	"github.com/guicostaarantes/psi-server/utils/hash"
	"github.com/guicostaarantes/psi-server/utils/identifier"
	"github.com/guicostaarantes/psi-server/utils/match"
	"github.com/guicostaarantes/psi-server/utils/orm"
	"github.com/guicostaarantes/psi-server/utils/serializing"
)

// CreateUserWithPasswordService is a service that creates users and assigns them a password immediately
type CreateUserWithPasswordService struct {
	HashUtil        hash.IHashUtil
	IdentifierUtil  identifier.IIdentifierUtil
	MatchUtil       match.IMatchUtil
	OrmUtil         orm.IOrmUtil
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

	result := s.OrmUtil.Db().Where("email = ?", userInput.Email).Limit(1).Find(&userWithSameEmail)
	if result.Error != nil {
		return result.Error
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
		ID:     id,
		Active: true,
		Email:  userInput.Email,
		Role:   userInput.Role,
	}

	user.Password = hashedPwd

	result = s.OrmUtil.Db().Create(&user)
	if result.Error != nil {
		return result.Error
	}

	return nil

}
