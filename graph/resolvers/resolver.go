package resolvers

import (
	characteristics_services "github.com/guicostaarantes/psi-server/modules/characteristics/services"
	mails_services "github.com/guicostaarantes/psi-server/modules/mails/services"
	profiles_services "github.com/guicostaarantes/psi-server/modules/profiles/services"
	schedule_services "github.com/guicostaarantes/psi-server/modules/schedule/services"
	translations_services "github.com/guicostaarantes/psi-server/modules/translations/services"
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
	MaxAffinityNumber                       int64
	SecondsLimitAvailability                int64
	SecondsMinimumAvailability              int64
	SecondsToCooldownReset                  int64
	SecondsToExpire                         int64
	SecondsToExpireReset                    int64
	askResetPasswordService                 *users_services.AskResetPasswordService
	assignTreatmentService                  *treatments_services.AssignTreatmentService
	authenticateUserService                 *users_services.AuthenticateUserService
	cancelAppointmentByPatientService       *schedule_services.CancelAppointmentByPatientService
	cancelAppointmentByPsychologistService  *schedule_services.CancelAppointmentByPsychologistService
	confirmAppointmentService               *schedule_services.ConfirmAppointmentService
	createTreatmentService                  *treatments_services.CreateTreatmentService
	createUserService                       *users_services.CreateUserService
	createUserWithPasswordService           *users_services.CreateUserWithPasswordService
	deleteTreatmentService                  *treatments_services.DeleteTreatmentService
	denyAppointmentService                  *schedule_services.DenyAppointmentService
	finalizeTreatmentService                *treatments_services.FinalizeTreatmentService
	getAppointmentsOfPatientService         *schedule_services.GetAppointmentsOfPatientService
	getAppointmentsOfPsychologistService    *schedule_services.GetAppointmentsOfPsychologistService
	getAvailabilityService                  *schedule_services.GetAvailabilityService
	getCharacteristicsByIDService           *characteristics_services.GetCharacteristicsByIDService
	getCharacteristicsService               *characteristics_services.GetCharacteristicsService
	getPatientByUserIDService               *profiles_services.GetPatientByUserIDService
	getPatientService                       *profiles_services.GetPatientService
	getPatientTreatmentsService             *treatments_services.GetPatientTreatmentsService
	getPreferencesByIDService               *characteristics_services.GetPreferencesByIDService
	getPsychologistByUserIDService          *profiles_services.GetPsychologistByUserIDService
	getPsychologistService                  *profiles_services.GetPsychologistService
	getPsychologistTreatmentsService        *treatments_services.GetPsychologistTreatmentsService
	getTopAffinitiesForPatientService       *characteristics_services.GetTopAffinitiesForPatientService
	getTranslationsService                  *translations_services.GetTranslationsService
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
	setTopAffinitiesForPatientService       *characteristics_services.SetTopAffinitiesForPatientService
	setTranslationsService                  *translations_services.SetTranslationsService
	updateTreatmentService                  *treatments_services.UpdateTreatmentService
	updateUserService                       *users_services.UpdateUserService
	upsertPatientService                    *profiles_services.UpsertPatientService
	upsertPsychologistService               *profiles_services.UpsertPsychologistService
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

// CancelAppointmentByPatientService gets or sets the service with same name
func (r *Resolver) CancelAppointmentByPatientService() *schedule_services.CancelAppointmentByPatientService {
	if r.cancelAppointmentByPatientService == nil {
		r.cancelAppointmentByPatientService = &schedule_services.CancelAppointmentByPatientService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.cancelAppointmentByPatientService
}

// CancelAppointmentByPsychologistService gets or sets the service with same name
func (r *Resolver) CancelAppointmentByPsychologistService() *schedule_services.CancelAppointmentByPsychologistService {
	if r.cancelAppointmentByPsychologistService == nil {
		r.cancelAppointmentByPsychologistService = &schedule_services.CancelAppointmentByPsychologistService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.cancelAppointmentByPsychologistService
}

// ConfirmAppointmentService gets or sets the service with same name
func (r *Resolver) ConfirmAppointmentService() *schedule_services.ConfirmAppointmentService {
	if r.confirmAppointmentService == nil {
		r.confirmAppointmentService = &schedule_services.ConfirmAppointmentService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.confirmAppointmentService
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

// DenyAppointmentService gets or sets the service with same name
func (r *Resolver) DenyAppointmentService() *schedule_services.DenyAppointmentService {
	if r.denyAppointmentService == nil {
		r.denyAppointmentService = &schedule_services.DenyAppointmentService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.denyAppointmentService
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

// GetTranslationsService gets or sets the service with same name
func (r *Resolver) GetTranslationsService() *translations_services.GetTranslationsService {
	if r.getTranslationsService == nil {
		r.getTranslationsService = &translations_services.GetTranslationsService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.getTranslationsService
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

// GetPsychologistService gets or sets the service with same name
func (r *Resolver) GetPsychologistService() *profiles_services.GetPsychologistService {
	if r.getPsychologistService == nil {
		r.getPsychologistService = &profiles_services.GetPsychologistService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.getPsychologistService
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

// GetPatientService gets or sets the service with same name
func (r *Resolver) GetPatientService() *profiles_services.GetPatientService {
	if r.getPatientService == nil {
		r.getPatientService = &profiles_services.GetPatientService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.getPatientService
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

// GetTopAffinitiesForPatientService gets or sets the service with same name
func (r *Resolver) GetTopAffinitiesForPatientService() *characteristics_services.GetTopAffinitiesForPatientService {
	if r.getTopAffinitiesForPatientService == nil {
		r.getTopAffinitiesForPatientService = &characteristics_services.GetTopAffinitiesForPatientService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.getTopAffinitiesForPatientService
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

// SetTranslationsService gets or sets the service with same name
func (r *Resolver) SetTranslationsService() *translations_services.SetTranslationsService {
	if r.setTranslationsService == nil {
		r.setTranslationsService = &translations_services.SetTranslationsService{
			DatabaseUtil: r.DatabaseUtil,
		}
	}
	return r.setTranslationsService
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

// SetTopAffinitiesForPatientService gets or sets the service with same name
func (r *Resolver) SetTopAffinitiesForPatientService() *characteristics_services.SetTopAffinitiesForPatientService {
	if r.setTopAffinitiesForPatientService == nil {
		r.setTopAffinitiesForPatientService = &characteristics_services.SetTopAffinitiesForPatientService{
			DatabaseUtil:      r.DatabaseUtil,
			MaxAffinityNumber: r.MaxAffinityNumber,
		}
	}
	return r.setTopAffinitiesForPatientService
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

// UpsertPatientService gets or sets the service with same name
func (r *Resolver) UpsertPatientService() *profiles_services.UpsertPatientService {
	if r.upsertPatientService == nil {
		r.upsertPatientService = &profiles_services.UpsertPatientService{
			DatabaseUtil:   r.DatabaseUtil,
			IdentifierUtil: r.IdentifierUtil,
		}
	}
	return r.upsertPatientService
}

// UpsertPsychologistService gets or sets the service with same name
func (r *Resolver) UpsertPsychologistService() *profiles_services.UpsertPsychologistService {
	if r.upsertPsychologistService == nil {
		r.upsertPsychologistService = &profiles_services.UpsertPsychologistService{
			DatabaseUtil:   r.DatabaseUtil,
			IdentifierUtil: r.IdentifierUtil,
		}
	}
	return r.upsertPsychologistService
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
