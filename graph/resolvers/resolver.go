package resolvers

import (
	mails_services "github.com/guicostaarantes/psi-server/modules/mails/services"
	users_services "github.com/guicostaarantes/psi-server/modules/users/services"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/hash"
	"github.com/guicostaarantes/psi-server/utils/identifier"
	"github.com/guicostaarantes/psi-server/utils/mail"
	"github.com/guicostaarantes/psi-server/utils/match"
	"github.com/guicostaarantes/psi-server/utils/serializing"
	"github.com/guicostaarantes/psi-server/utils/token"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DatabaseUtil                  database.IDatabaseUtil
	HashUtil                      hash.IHashUtil
	IdentifierUtil                identifier.IIdentifierUtil
	MailUtil                      mail.IMailUtil
	MatchUtil                     match.IMatchUtil
	SerializingUtil               serializing.ISerializingUtil
	TokenUtil                     token.ITokenUtil
	SecondsToExpire               int64
	SecondsToExpireReset          int64
	activateUserService           *users_services.ActivateUserService
	authenticateUserService       *users_services.AuthenticateUserService
	createUserService             *users_services.CreateUserService
	createUserWithPasswordService *users_services.CreateUserWithPasswordService
	getUserByIdService            *users_services.GetUserByIdService
	processPendingMailsService    *mails_services.ProcessPendingMailsService
	resetPasswordService          *users_services.ResetPasswordService
	validateUserTokenService      *users_services.ValidateUserTokenService
}

func (r *Resolver) ActivateUserService() *users_services.ActivateUserService {
	if r.activateUserService == nil {
		r.activateUserService = &users_services.ActivateUserService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.activateUserService
}

func (r *Resolver) AuthenticateUserService() *users_services.AuthenticateUserService {
	if r.authenticateUserService == nil {
		r.authenticateUserService = &users_services.AuthenticateUserService{
			DatabaseUtil:    r.DatabaseUtil,
			HashUtil:        r.HashUtil,
			IdentifierUtil:  r.IdentifierUtil,
			SerializingUtil: r.SerializingUtil,
			TokenUtil:       r.TokenUtil,
			SecondsToExpire: r.SecondsToExpire,
		}
	}
	return r.authenticateUserService
}

func (r *Resolver) CreateUserService() *users_services.CreateUserService {
	if r.createUserService == nil {
		r.createUserService = &users_services.CreateUserService{
			DatabaseUtil:    r.DatabaseUtil,
			IdentifierUtil:  r.IdentifierUtil,
			MatchUtil:       r.MatchUtil,
			SerializingUtil: r.SerializingUtil,
			TokenUtil:       r.TokenUtil,
			SecondsToExpire: r.SecondsToExpire,
		}
	}
	return r.createUserService
}

func (r *Resolver) CreateUserWithPasswordService() *users_services.CreateUserWithPasswordService {
	if r.createUserWithPasswordService == nil {
		r.createUserWithPasswordService = &users_services.CreateUserWithPasswordService{
			DatabaseUtil:    r.DatabaseUtil,
			HashUtil:        r.HashUtil,
			IdentifierUtil:  r.IdentifierUtil,
			MatchUtil:       r.MatchUtil,
			SerializingUtil: r.SerializingUtil,
		}
	}
	return r.createUserWithPasswordService
}

func (r *Resolver) GetUserByIdService() *users_services.GetUserByIdService {
	if r.getUserByIdService == nil {
		r.getUserByIdService = &users_services.GetUserByIdService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.getUserByIdService
}

func (r *Resolver) ProcessPendingMailsService() *mails_services.ProcessPendingMailsService {
	if r.processPendingMailsService == nil {
		r.processPendingMailsService = &mails_services.ProcessPendingMailsService{
			DatabaseUtil: r.DatabaseUtil,
			MailUtil:     r.MailUtil,
		}
	}
	return r.processPendingMailsService
}

func (r *Resolver) ResetPasswordService() *users_services.ResetPasswordService {
	if r.resetPasswordService == nil {
		r.resetPasswordService = &users_services.ResetPasswordService{
			DatabaseUtil: r.DatabaseUtil,
			HashUtil:     r.HashUtil,
		}
	}
	return r.resetPasswordService
}

func (r *Resolver) ValidateUserTokenService() *users_services.ValidateUserTokenService {
	if r.validateUserTokenService == nil {
		r.validateUserTokenService = &users_services.ValidateUserTokenService{
			DatabaseUtil:    r.DatabaseUtil,
			SerializingUtil: r.SerializingUtil,
			SecondsToExpire: r.SecondsToExpire,
		}
	}
	return r.validateUserTokenService
}
