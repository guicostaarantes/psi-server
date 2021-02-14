package resolvers

import (
	characteristics_services "github.com/guicostaarantes/psi-server/modules/characteristics/services"
	mails_services "github.com/guicostaarantes/psi-server/modules/mails/services"
	profiles_services "github.com/guicostaarantes/psi-server/modules/profiles/services"
	schedule_services "github.com/guicostaarantes/psi-server/modules/schedule/services"
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
	DatabaseUtil                       database.IDatabaseUtil
	HashUtil                           hash.IHashUtil
	IdentifierUtil                     identifier.IIdentifierUtil
	MailUtil                           mail.IMailUtil
	MatchUtil                          match.IMatchUtil
	SerializingUtil                    serializing.ISerializingUtil
	TokenUtil                          token.ITokenUtil
	SecondsLimitAvailability           int64
	SecondsMinimumAvailability         int64
	SecondsToCooldownReset             int64
	SecondsToExpire                    int64
	SecondsToExpireReset               int64
	askResetPasswordService            *users_services.AskResetPasswordService
	assignSlotService                  *schedule_services.AssignSlotService
	authenticateUserService            *users_services.AuthenticateUserService
	createPatientService               *profiles_services.CreatePatientService
	createPsychologistService          *profiles_services.CreatePsychologistService
	createSlotService                  *schedule_services.CreateSlotService
	createUserService                  *users_services.CreateUserService
	createUserWithPasswordService      *users_services.CreateUserWithPasswordService
	deleteSlotService                  *schedule_services.DeleteSlotService
	finalizeSlotService                *schedule_services.FinalizeSlotService
	getAvailabilityService             *schedule_services.GetAvailabilityService
	getCharacteristicsByIDService      *characteristics_services.GetCharacteristicsByIDService
	getCharacteristicsService          *characteristics_services.GetCharacteristicsService
	getPatientByUserIDService          *profiles_services.GetPatientByUserIDService
	getPatientSlotsService             *schedule_services.GetPatientSlotsService
	getPreferencesByIDService          *characteristics_services.GetPreferencesByIDService
	getPsychologistByUserIDService     *profiles_services.GetPsychologistByUserIDService
	getPsychologistSlotsService        *schedule_services.GetPsychologistSlotsService
	getUserByIDService                 *users_services.GetUserByIDService
	getUsersByRoleService              *users_services.GetUsersByRoleService
	interruptSlotByPatientService      *schedule_services.InterruptSlotByPatientService
	interruptSlotByPsychologistService *schedule_services.InterruptSlotByPsychologistService
	processPendingMailsService         *mails_services.ProcessPendingMailsService
	resetPasswordService               *users_services.ResetPasswordService
	setAvailabilityService             *schedule_services.SetAvailabilityService
	setCharacteristicChoicesService    *characteristics_services.SetCharacteristicChoicesService
	setCharacteristicsService          *characteristics_services.SetCharacteristicsService
	setPreferencesService              *characteristics_services.SetPreferencesService
	updatePatientService               *profiles_services.UpdatePatientService
	updatePsychologistService          *profiles_services.UpdatePsychologistService
	updateSlotService                  *schedule_services.UpdateSlotService
	updateUserService                  *users_services.UpdateUserService
	validateUserTokenService           *users_services.ValidateUserTokenService
}

// AskResetPasswordService gets or sets the service with same name
func (r *Resolver) AskResetPasswordService() *users_services.AskResetPasswordService {
	if r.askResetPasswordService == nil {
		r.askResetPasswordService = &users_services.AskResetPasswordService{
			DatabaseUtil:    r.DatabaseUtil,
			IdentifierUtil:  r.IdentifierUtil,
			TokenUtil:       r.TokenUtil,
			SecondsToExpire: r.SecondsToExpire,
		}
	}
	return r.askResetPasswordService
}

// AssignSlotService gets or sets the service with same name
func (r *Resolver) AssignSlotService() *schedule_services.AssignSlotService {
	if r.assignSlotService == nil {
		r.assignSlotService = &schedule_services.AssignSlotService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.assignSlotService
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

// CreateSlotService gets or sets the service with same name
func (r *Resolver) CreateSlotService() *schedule_services.CreateSlotService {
	if r.createSlotService == nil {
		r.createSlotService = &schedule_services.CreateSlotService{
			DatabaseUtil:   r.DatabaseUtil,
			IdentifierUtil: r.IdentifierUtil,
		}
	}
	return r.createSlotService
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

// DeleteSlotService gets or sets the service with same name
func (r *Resolver) DeleteSlotService() *schedule_services.DeleteSlotService {
	if r.deleteSlotService == nil {
		r.deleteSlotService = &schedule_services.DeleteSlotService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.deleteSlotService
}

// FinalizeSlotService gets or sets the service with same name
func (r *Resolver) FinalizeSlotService() *schedule_services.FinalizeSlotService {
	if r.finalizeSlotService == nil {
		r.finalizeSlotService = &schedule_services.FinalizeSlotService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.finalizeSlotService
}

// GetAvailabilityService gets or sets the service with same name
func (r *Resolver) GetAvailabilityService() *schedule_services.GetAvailabilityService {
	if r.getAvailabilityService == nil {
		r.getAvailabilityService = &schedule_services.GetAvailabilityService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.getAvailabilityService
}

// GetCharacteristicsByIDService gets or sets the service with same name
func (r *Resolver) GetCharacteristicsByIDService() *characteristics_services.GetCharacteristicsByIDService {
	if r.getCharacteristicsByIDService == nil {
		r.getCharacteristicsByIDService = &characteristics_services.GetCharacteristicsByIDService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.getCharacteristicsByIDService
}

// GetCharacteristicsService gets or sets the service with same name
func (r *Resolver) GetCharacteristicsService() *characteristics_services.GetCharacteristicsService {
	if r.getCharacteristicsService == nil {
		r.getCharacteristicsService = &characteristics_services.GetCharacteristicsService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.getCharacteristicsService
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

// GetPatientSlotsService gets or sets the service with same name
func (r *Resolver) GetPatientSlotsService() *schedule_services.GetPatientSlotsService {
	if r.getPatientSlotsService == nil {
		r.getPatientSlotsService = &schedule_services.GetPatientSlotsService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.getPatientSlotsService
}

// GetPreferencesByIDService gets or sets the service with same name
func (r *Resolver) GetPreferencesByIDService() *characteristics_services.GetPreferencesByIDService {
	if r.getPreferencesByIDService == nil {
		r.getPreferencesByIDService = &characteristics_services.GetPreferencesByIDService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.getPreferencesByIDService
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

// GetPsychologistSlotsService gets or sets the service with same name
func (r *Resolver) GetPsychologistSlotsService() *schedule_services.GetPsychologistSlotsService {
	if r.getPsychologistSlotsService == nil {
		r.getPsychologistSlotsService = &schedule_services.GetPsychologistSlotsService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.getPsychologistSlotsService
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

// InterruptSlotByPatientService gets or sets the service with same name
func (r *Resolver) InterruptSlotByPatientService() *schedule_services.InterruptSlotByPatientService {
	if r.interruptSlotByPatientService == nil {
		r.interruptSlotByPatientService = &schedule_services.InterruptSlotByPatientService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.interruptSlotByPatientService
}

// InterruptSlotByPsychologistService gets or sets the service with same name
func (r *Resolver) InterruptSlotByPsychologistService() *schedule_services.InterruptSlotByPsychologistService {
	if r.interruptSlotByPsychologistService == nil {
		r.interruptSlotByPsychologistService = &schedule_services.InterruptSlotByPsychologistService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.interruptSlotByPsychologistService
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

// SetAvailabilityService gets or sets the service with same name
func (r *Resolver) SetAvailabilityService() *schedule_services.SetAvailabilityService {
	if r.setAvailabilityService == nil {
		r.setAvailabilityService = &schedule_services.SetAvailabilityService{
			DatabaseUtil:               r.DatabaseUtil,
			SecondsLimitAvailability:   r.SecondsLimitAvailability,
			SecondsMinimumAvailability: r.SecondsMinimumAvailability,
		}
	}
	return r.setAvailabilityService
}

// SetCharacteristicChoicesService gets or sets the service with same name
func (r *Resolver) SetCharacteristicChoicesService() *characteristics_services.SetCharacteristicChoicesService {
	if r.setCharacteristicChoicesService == nil {
		r.setCharacteristicChoicesService = &characteristics_services.SetCharacteristicChoicesService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.setCharacteristicChoicesService
}

// SetCharacteristicsService gets or sets the service with same name
func (r *Resolver) SetCharacteristicsService() *characteristics_services.SetCharacteristicsService {
	if r.setCharacteristicsService == nil {
		r.setCharacteristicsService = &characteristics_services.SetCharacteristicsService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.setCharacteristicsService
}

// SetPreferencesService gets or sets the service with same name
func (r *Resolver) SetPreferencesService() *characteristics_services.SetPreferencesService {
	if r.setPreferencesService == nil {
		r.setPreferencesService = &characteristics_services.SetPreferencesService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.setPreferencesService
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

// UpdateSlotService gets or sets the service with same name
func (r *Resolver) UpdateSlotService() *schedule_services.UpdateSlotService {
	if r.updateSlotService == nil {
		r.updateSlotService = &schedule_services.UpdateSlotService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.updateSlotService
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
