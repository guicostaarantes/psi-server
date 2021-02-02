package resolvers

import (
	mails_services "github.com/guicostaarantes/psi-server/modules/mails/services"
	profiles_services "github.com/guicostaarantes/psi-server/modules/profiles/services"
	users_services "github.com/guicostaarantes/psi-server/modules/users/services"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/hash"
	"github.com/guicostaarantes/psi-server/utils/identifier"
	"github.com/guicostaarantes/psi-server/utils/mail"
	"github.com/guicostaarantes/psi-server/utils/match"
	"github.com/guicostaarantes/psi-server/utils/merge"
	"github.com/guicostaarantes/psi-server/utils/serializing"
	"github.com/guicostaarantes/psi-server/utils/token"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DatabaseUtil                   database.IDatabaseUtil
	HashUtil                       hash.IHashUtil
	IdentifierUtil                 identifier.IIdentifierUtil
	MailUtil                       mail.IMailUtil
	MatchUtil                      match.IMatchUtil
	MergeUtil                      merge.IMergeUtil
	SerializingUtil                serializing.ISerializingUtil
	TokenUtil                      token.ITokenUtil
	SecondsToCooldownReset         int64
	SecondsToExpire                int64
	SecondsToExpireReset           int64
	activateUserService            *users_services.ActivateUserService
	askResetPasswordService        *users_services.AskResetPasswordService
	authenticateUserService        *users_services.AuthenticateUserService
	createPsyCharacteristicService *profiles_services.CreatePsyCharacteristicService
	createPsychologistService      *profiles_services.CreatePsychologistService
	createUserService              *users_services.CreateUserService
	createUserWithPasswordService  *users_services.CreateUserWithPasswordService
	getPsychologistByUserIDService *profiles_services.GetPsychologistByUserIDService
	getUsersByRoleService          *users_services.GetUsersByRoleService
	getUserByIdService             *users_services.GetUserByIdService
	processPendingMailsService     *mails_services.ProcessPendingMailsService
	resetPasswordService           *users_services.ResetPasswordService
	updatePsyCharacteristicService *profiles_services.UpdatePsyCharacteristicService
	updatePsychologistService      *profiles_services.UpdatePsychologistService
	updateUserService              *users_services.UpdateUserService
	validateUserTokenService       *users_services.ValidateUserTokenService
}

func (r *Resolver) ActivateUserService() *users_services.ActivateUserService {
	if r.activateUserService == nil {
		r.activateUserService = &users_services.ActivateUserService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.activateUserService
}

func (r *Resolver) AskResetPasswordService() *users_services.AskResetPasswordService {
	if r.askResetPasswordService == nil {
		r.askResetPasswordService = &users_services.AskResetPasswordService{
			DatabaseUtil:      r.DatabaseUtil,
			IdentifierUtil:    r.IdentifierUtil,
			TokenUtil:         r.TokenUtil,
			SecondsToCooldown: r.SecondsToCooldownReset,
			SecondsToExpire:   r.SecondsToExpire,
		}
	}
	return r.askResetPasswordService
}

func (r *Resolver) AuthenticateUserService() *users_services.AuthenticateUserService {
	if r.authenticateUserService == nil {
		r.authenticateUserService = &users_services.AuthenticateUserService{
			DatabaseUtil:    r.DatabaseUtil,
			HashUtil:        r.HashUtil,
			SerializingUtil: r.SerializingUtil,
			TokenUtil:       r.TokenUtil,
			SecondsToExpire: r.SecondsToExpire,
		}
	}
	return r.authenticateUserService
}

func (r *Resolver) CreatePsyCharacteristicService() *profiles_services.CreatePsyCharacteristicService {
	if r.createPsyCharacteristicService == nil {
		r.createPsyCharacteristicService = &profiles_services.CreatePsyCharacteristicService{
			DatabaseUtil:   r.DatabaseUtil,
			IdentifierUtil: r.IdentifierUtil,
			MergeUtil:      r.MergeUtil,
		}
	}
	return r.createPsyCharacteristicService
}

func (r *Resolver) CreatePsychologistService() *profiles_services.CreatePsychologistService {
	if r.createPsychologistService == nil {
		r.createPsychologistService = &profiles_services.CreatePsychologistService{
			DatabaseUtil:   r.DatabaseUtil,
			IdentifierUtil: r.IdentifierUtil,
			MergeUtil:      r.MergeUtil,
		}
	}
	return r.createPsychologistService
}

func (r *Resolver) CreateUserService() *users_services.CreateUserService {
	if r.createUserService == nil {
		r.createUserService = &users_services.CreateUserService{
			DatabaseUtil:    r.DatabaseUtil,
			IdentifierUtil:  r.IdentifierUtil,
			MatchUtil:       r.MatchUtil,
			MergeUtil:       r.MergeUtil,
			SerializingUtil: r.SerializingUtil,
			TokenUtil:       r.TokenUtil,
			SecondsToExpire: r.SecondsToExpireReset,
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
			MergeUtil:       r.MergeUtil,
			SerializingUtil: r.SerializingUtil,
		}
	}
	return r.createUserWithPasswordService
}

func (r *Resolver) GetPsychologistByUserIDService() *profiles_services.GetPsychologistByUserIDService {
	if r.getPsychologistByUserIDService == nil {
		r.getPsychologistByUserIDService = &profiles_services.GetPsychologistByUserIDService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.getPsychologistByUserIDService
}

func (r *Resolver) GetUsersByRoleService() *users_services.GetUsersByRoleService {
	if r.getUsersByRoleService == nil {
		r.getUsersByRoleService = &users_services.GetUsersByRoleService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.getUsersByRoleService
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
			MatchUtil:    r.MatchUtil,
		}
	}
	return r.resetPasswordService
}

func (r *Resolver) UpdatePsyCharacteristicService() *profiles_services.UpdatePsyCharacteristicService {
	if r.updatePsyCharacteristicService == nil {
		r.updatePsyCharacteristicService = &profiles_services.UpdatePsyCharacteristicService{
			DatabaseUtil:   r.DatabaseUtil,
			IdentifierUtil: r.IdentifierUtil,
			MergeUtil:      r.MergeUtil,
		}
	}
	return r.updatePsyCharacteristicService
}

func (r *Resolver) UpdatePsychologistService() *profiles_services.UpdatePsychologistService {
	if r.updatePsychologistService == nil {
		r.updatePsychologistService = &profiles_services.UpdatePsychologistService{
			DatabaseUtil: r.DatabaseUtil,
			MergeUtil:    r.MergeUtil,
		}
	}
	return r.updatePsychologistService
}

func (r *Resolver) UpdateUserService() *users_services.UpdateUserService {
	if r.updateUserService == nil {
		r.updateUserService = &users_services.UpdateUserService{
			DatabaseUtil: r.DatabaseUtil,
			MergeUtil:    r.MergeUtil,
		}
	}
	return r.updateUserService
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
