package resolvers

import (
	characteristics_services "github.com/guicostaarantes/psi-server/modules/characteristics/services"
	mails_services "github.com/guicostaarantes/psi-server/modules/mails/services"
	profiles_services "github.com/guicostaarantes/psi-server/modules/profiles/services"
	schedule_services "github.com/guicostaarantes/psi-server/modules/schedule/services"
	treatments_services "github.com/guicostaarantes/psi-server/modules/treatments/services"
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
	DatabaseUtil                            database.IDatabaseUtil
	HashUtil                                hash.IHashUtil
	IdentifierUtil                          identifier.IIdentifierUtil
	MailUtil                                mail.IMailUtil
	MatchUtil                               match.IMatchUtil
	SerializingUtil                         serializing.ISerializingUtil
	TokenUtil                               token.ITokenUtil
	SecondsLimitAvailability                int64
	SecondsMinimumAvailability              int64
	SecondsToCooldownReset                  int64
	SecondsToExpire                         int64
	SecondsToExpireReset                    int64
	askResetPasswordService                 *users_services.AskResetPasswordService
	assignTreatmentService                  *treatments_services.AssignTreatmentService
	authenticateUserService                 *users_services.AuthenticateUserService
	createPatientService                    *profiles_services.CreatePatientService
	createPsychologistService               *profiles_services.CreatePsychologistService
	createTreatmentService                  *treatments_services.CreateTreatmentService
	createUserService                       *users_services.CreateUserService
	createUserWithPasswordService           *users_services.CreateUserWithPasswordService
	deleteTreatmentService                  *treatments_services.DeleteTreatmentService
	finalizeTreatmentService                *treatments_services.FinalizeTreatmentService
	getAppointmentsOfPatientService         *schedule_services.GetAppointmentsOfPatientService
	getAppointmentsOfPsychologistService    *schedule_services.GetAppointmentsOfPsychologistService
	getAvailabilityService                  *schedule_services.GetAvailabilityService
	getCharacteristicsByIDService           *characteristics_services.GetCharacteristicsByIDService
	getCharacteristicsService               *characteristics_services.GetCharacteristicsService
	getPatientByUserIDService               *profiles_services.GetPatientByUserIDService
	getPatientTreatmentsService             *treatments_services.GetPatientTreatmentsService
	getPreferencesByIDService               *characteristics_services.GetPreferencesByIDService
	getPsychologistByUserIDService          *profiles_services.GetPsychologistByUserIDService
	getPsychologistTreatmentsService        *treatments_services.GetPsychologistTreatmentsService
	getUserByIDService                      *users_services.GetUserByIDService
	getUsersByRoleService                   *users_services.GetUsersByRoleService
	interruptTreatmentByPatientService      *treatments_services.InterruptTreatmentByPatientService
	interruptTreatmentByPsychologistService *treatments_services.InterruptTreatmentByPsychologistService
	processPendingMailsService              *mails_services.ProcessPendingMailsService
	proposeAppointmentService               *schedule_services.ProposeAppointmentService
	resetPasswordService                    *users_services.ResetPasswordService
	setAvailabilityService                  *schedule_services.SetAvailabilityService
	setCharacteristicChoicesService         *characteristics_services.SetCharacteristicChoicesService
	setCharacteristicsService               *characteristics_services.SetCharacteristicsService
	setPreferencesService                   *characteristics_services.SetPreferencesService
	updatePatientService                    *profiles_services.UpdatePatientService
	updatePsychologistService               *profiles_services.UpdatePsychologistService
	updateTreatmentService                  *treatments_services.UpdateTreatmentService
	updateUserService                       *users_services.UpdateUserService
	validateUserTokenService                *users_services.ValidateUserTokenService
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

// AssignTreatmentService gets or sets the service with same name
func (r *Resolver) AssignTreatmentService() *treatments_services.AssignTreatmentService {
	if r.assignTreatmentService == nil {
		r.assignTreatmentService = &treatments_services.AssignTreatmentService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.assignTreatmentService
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

// CreateTreatmentService gets or sets the service with same name
func (r *Resolver) CreateTreatmentService() *treatments_services.CreateTreatmentService {
	if r.createTreatmentService == nil {
		r.createTreatmentService = &treatments_services.CreateTreatmentService{
			DatabaseUtil:   r.DatabaseUtil,
			IdentifierUtil: r.IdentifierUtil,
		}
	}
	return r.createTreatmentService
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

// DeleteTreatmentService gets or sets the service with same name
func (r *Resolver) DeleteTreatmentService() *treatments_services.DeleteTreatmentService {
	if r.deleteTreatmentService == nil {
		r.deleteTreatmentService = &treatments_services.DeleteTreatmentService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.deleteTreatmentService
}

// FinalizeTreatmentService gets or sets the service with same name
func (r *Resolver) FinalizeTreatmentService() *treatments_services.FinalizeTreatmentService {
	if r.finalizeTreatmentService == nil {
		r.finalizeTreatmentService = &treatments_services.FinalizeTreatmentService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.finalizeTreatmentService
}

// GetAppointmentsOfPatientService gets or sets the service with same name
func (r *Resolver) GetAppointmentsOfPatientService() *schedule_services.GetAppointmentsOfPatientService {
	if r.getAppointmentsOfPatientService == nil {
		r.getAppointmentsOfPatientService = &schedule_services.GetAppointmentsOfPatientService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.getAppointmentsOfPatientService
}

// GetAppointmentsOfPsychologistService gets or sets the service with same name
func (r *Resolver) GetAppointmentsOfPsychologistService() *schedule_services.GetAppointmentsOfPsychologistService {
	if r.getAppointmentsOfPsychologistService == nil {
		r.getAppointmentsOfPsychologistService = &schedule_services.GetAppointmentsOfPsychologistService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.getAppointmentsOfPsychologistService
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

// GetPatientTreatmentsService gets or sets the service with same name
func (r *Resolver) GetPatientTreatmentsService() *treatments_services.GetPatientTreatmentsService {
	if r.getPatientTreatmentsService == nil {
		r.getPatientTreatmentsService = &treatments_services.GetPatientTreatmentsService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.getPatientTreatmentsService
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

// GetPsychologistTreatmentsService gets or sets the service with same name
func (r *Resolver) GetPsychologistTreatmentsService() *treatments_services.GetPsychologistTreatmentsService {
	if r.getPsychologistTreatmentsService == nil {
		r.getPsychologistTreatmentsService = &treatments_services.GetPsychologistTreatmentsService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.getPsychologistTreatmentsService
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

// InterruptTreatmentByPatientService gets or sets the service with same name
func (r *Resolver) InterruptTreatmentByPatientService() *treatments_services.InterruptTreatmentByPatientService {
	if r.interruptTreatmentByPatientService == nil {
		r.interruptTreatmentByPatientService = &treatments_services.InterruptTreatmentByPatientService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.interruptTreatmentByPatientService
}

// InterruptTreatmentByPsychologistService gets or sets the service with same name
func (r *Resolver) InterruptTreatmentByPsychologistService() *treatments_services.InterruptTreatmentByPsychologistService {
	if r.interruptTreatmentByPsychologistService == nil {
		r.interruptTreatmentByPsychologistService = &treatments_services.InterruptTreatmentByPsychologistService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.interruptTreatmentByPsychologistService
}

// ProposeAppointmentService gets or sets the service with same name
func (r *Resolver) ProposeAppointmentService() *schedule_services.ProposeAppointmentService {
	if r.proposeAppointmentService == nil {
		r.proposeAppointmentService = &schedule_services.ProposeAppointmentService{
			DatabaseUtil:   r.DatabaseUtil,
			IdentifierUtil: r.IdentifierUtil,
		}
	}
	return r.proposeAppointmentService
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

// UpdateTreatmentService gets or sets the service with same name
func (r *Resolver) UpdateTreatmentService() *treatments_services.UpdateTreatmentService {
	if r.updateTreatmentService == nil {
		r.updateTreatmentService = &treatments_services.UpdateTreatmentService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.updateTreatmentService
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
