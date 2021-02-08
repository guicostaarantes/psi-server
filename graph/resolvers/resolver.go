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
	"github.com/guicostaarantes/psi-server/utils/serializing"
	"github.com/guicostaarantes/psi-server/utils/token"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver receives all utils and registers all services within the application
type Resolver struct {
	DatabaseUtil                                 database.IDatabaseUtil
	HashUtil                                     hash.IHashUtil
	IdentifierUtil                               identifier.IIdentifierUtil
	MailUtil                                     mail.IMailUtil
	MatchUtil                                    match.IMatchUtil
	SerializingUtil                              serializing.ISerializingUtil
	TokenUtil                                    token.ITokenUtil
	SecondsToCooldownReset                       int64
	SecondsToExpire                              int64
	SecondsToExpireReset                         int64
	activateUserService                          *users_services.ActivateUserService
	askResetPasswordService                      *users_services.AskResetPasswordService
	authenticateUserService                      *users_services.AuthenticateUserService
	createPatientService                         *profiles_services.CreatePatientService
	createPsychologistService                    *profiles_services.CreatePsychologistService
	createUserService                            *users_services.CreateUserService
	createUserWithPasswordService                *users_services.CreateUserWithPasswordService
	getPatientByUserIDService                    *profiles_services.GetPatientByUserIDService
	getPatientCharacteristicsService             *profiles_services.GetPatientCharacteristicsService
	getPsychologistCharacteristicsByPsyIDService *profiles_services.GetPsychologistCharacteristicsByPsyIDService
	getPsychologistCharacteristicsService        *profiles_services.GetPsychologistCharacteristicsService
	getPsychologistByUserIDService               *profiles_services.GetPsychologistByUserIDService
	getUsersByRoleService                        *users_services.GetUsersByRoleService
	getUserByIDService                           *users_services.GetUserByIDService
	processPendingMailsService                   *mails_services.ProcessPendingMailsService
	resetPasswordService                         *users_services.ResetPasswordService
	setPatientCharacteristicsService             *profiles_services.SetPatientCharacteristicsService
	setPsychologistCharacteristicChoicesService  *profiles_services.SetPsychologistCharacteristicChoicesService
	setPsychologistCharacteristicsService        *profiles_services.SetPsychologistCharacteristicsService
	updatePatientService                         *profiles_services.UpdatePatientService
	updatePsychologistService                    *profiles_services.UpdatePsychologistService
	updateUserService                            *users_services.UpdateUserService
	validateUserTokenService                     *users_services.ValidateUserTokenService
}

// ActivateUserService gets or sets the service with same name
func (r *Resolver) ActivateUserService() *users_services.ActivateUserService {
	if r.activateUserService == nil {
		r.activateUserService = &users_services.ActivateUserService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.activateUserService
}

// AskResetPasswordService gets or sets the service with same name
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

// AuthenticateUserService gets or sets the service with same name
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

// CreatePatientService gets or sets the service with same name
func (r *Resolver) CreatePatientService() *profiles_services.CreatePatientService {
	if r.createPatientService == nil {
		r.createPatientService = &profiles_services.CreatePatientService{
			DatabaseUtil:   r.DatabaseUtil,
			IdentifierUtil: r.IdentifierUtil,
		}
	}
	return r.createPatientService
}

// CreatePsychologistService gets or sets the service with same name
func (r *Resolver) CreatePsychologistService() *profiles_services.CreatePsychologistService {
	if r.createPsychologistService == nil {
		r.createPsychologistService = &profiles_services.CreatePsychologistService{
			DatabaseUtil:   r.DatabaseUtil,
			IdentifierUtil: r.IdentifierUtil,
		}
	}
	return r.createPsychologistService
}

// CreateUserService gets or sets the service with same name
func (r *Resolver) CreateUserService() *users_services.CreateUserService {
	if r.createUserService == nil {
		r.createUserService = &users_services.CreateUserService{
			DatabaseUtil:    r.DatabaseUtil,
			IdentifierUtil:  r.IdentifierUtil,
			MatchUtil:       r.MatchUtil,
			SerializingUtil: r.SerializingUtil,
			TokenUtil:       r.TokenUtil,
			SecondsToExpire: r.SecondsToExpireReset,
		}
	}
	return r.createUserService
}

// CreateUserWithPasswordService gets or sets the service with same name
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

// GetPatientByUserIDService gets or sets the service with same name
func (r *Resolver) GetPatientByUserIDService() *profiles_services.GetPatientByUserIDService {
	if r.getPatientByUserIDService == nil {
		r.getPatientByUserIDService = &profiles_services.GetPatientByUserIDService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.getPatientByUserIDService
}

// GetPatientCharacteristicsService gets or sets the service with same name
func (r *Resolver) GetPatientCharacteristicsService() *profiles_services.GetPatientCharacteristicsService {
	if r.getPatientCharacteristicsService == nil {
		r.getPatientCharacteristicsService = &profiles_services.GetPatientCharacteristicsService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.getPatientCharacteristicsService
}

// GetPsychologistCharacteristicsByPsyIDService gets or sets the service with same name
func (r *Resolver) GetPsychologistCharacteristicsByPsyIDService() *profiles_services.GetPsychologistCharacteristicsByPsyIDService {
	if r.getPsychologistCharacteristicsByPsyIDService == nil {
		r.getPsychologistCharacteristicsByPsyIDService = &profiles_services.GetPsychologistCharacteristicsByPsyIDService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.getPsychologistCharacteristicsByPsyIDService
}

// GetPsychologistCharacteristicsService gets or sets the service with same name
func (r *Resolver) GetPsychologistCharacteristicsService() *profiles_services.GetPsychologistCharacteristicsService {
	if r.getPsychologistCharacteristicsService == nil {
		r.getPsychologistCharacteristicsService = &profiles_services.GetPsychologistCharacteristicsService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.getPsychologistCharacteristicsService
}

// GetPsychologistByUserIDService gets or sets the service with same name
func (r *Resolver) GetPsychologistByUserIDService() *profiles_services.GetPsychologistByUserIDService {
	if r.getPsychologistByUserIDService == nil {
		r.getPsychologistByUserIDService = &profiles_services.GetPsychologistByUserIDService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.getPsychologistByUserIDService
}

// GetUsersByRoleService gets or sets the service with same name
func (r *Resolver) GetUsersByRoleService() *users_services.GetUsersByRoleService {
	if r.getUsersByRoleService == nil {
		r.getUsersByRoleService = &users_services.GetUsersByRoleService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.getUsersByRoleService
}

// GetUserByIDService gets or sets the service with same name
func (r *Resolver) GetUserByIDService() *users_services.GetUserByIDService {
	if r.getUserByIDService == nil {
		r.getUserByIDService = &users_services.GetUserByIDService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.getUserByIDService
}

// ProcessPendingMailsService gets or sets the service with same name
func (r *Resolver) ProcessPendingMailsService() *mails_services.ProcessPendingMailsService {
	if r.processPendingMailsService == nil {
		r.processPendingMailsService = &mails_services.ProcessPendingMailsService{
			DatabaseUtil: r.DatabaseUtil,
			MailUtil:     r.MailUtil,
		}
	}
	return r.processPendingMailsService
}

// SetPatientCharacteristicsService gets or sets the service with same name
func (r *Resolver) SetPatientCharacteristicsService() *profiles_services.SetPatientCharacteristicsService {
	if r.setPatientCharacteristicsService == nil {
		r.setPatientCharacteristicsService = &profiles_services.SetPatientCharacteristicsService{
			DatabaseUtil:   r.DatabaseUtil,
			IdentifierUtil: r.IdentifierUtil,
		}
	}
	return r.setPatientCharacteristicsService
}

// SetPsychologistCharacteristicChoicesService gets or sets the service with same name
func (r *Resolver) SetPsychologistCharacteristicChoicesService() *profiles_services.SetPsychologistCharacteristicChoicesService {
	if r.setPsychologistCharacteristicChoicesService == nil {
		r.setPsychologistCharacteristicChoicesService = &profiles_services.SetPsychologistCharacteristicChoicesService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.setPsychologistCharacteristicChoicesService
}

// SetPsychologistCharacteristicsService gets or sets the service with same name
func (r *Resolver) SetPsychologistCharacteristicsService() *profiles_services.SetPsychologistCharacteristicsService {
	if r.setPsychologistCharacteristicsService == nil {
		r.setPsychologistCharacteristicsService = &profiles_services.SetPsychologistCharacteristicsService{
			DatabaseUtil:   r.DatabaseUtil,
			IdentifierUtil: r.IdentifierUtil,
		}
	}
	return r.setPsychologistCharacteristicsService
}

// ResetPasswordService gets or sets the service with same name
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

// UpdatePatientService gets or sets the service with same name
func (r *Resolver) UpdatePatientService() *profiles_services.UpdatePatientService {
	if r.updatePatientService == nil {
		r.updatePatientService = &profiles_services.UpdatePatientService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.updatePatientService
}

// UpdatePsychologistService gets or sets the service with same name
func (r *Resolver) UpdatePsychologistService() *profiles_services.UpdatePsychologistService {
	if r.updatePsychologistService == nil {
		r.updatePsychologistService = &profiles_services.UpdatePsychologistService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.updatePsychologistService
}

// UpdateUserService gets or sets the service with same name
func (r *Resolver) UpdateUserService() *users_services.UpdateUserService {
	if r.updateUserService == nil {
		r.updateUserService = &users_services.UpdateUserService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.updateUserService
}

// ValidateUserTokenService gets or sets the service with same name
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
