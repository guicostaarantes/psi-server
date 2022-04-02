package resolvers

import (
	"time"

	agreements_services "github.com/guicostaarantes/psi-server/modules/agreements/services"
	appointments_services "github.com/guicostaarantes/psi-server/modules/appointments/services"
	characteristics_services "github.com/guicostaarantes/psi-server/modules/characteristics/services"
	cooldowns_services "github.com/guicostaarantes/psi-server/modules/cooldowns/services"
	files_services "github.com/guicostaarantes/psi-server/modules/files/services"
	mails_services "github.com/guicostaarantes/psi-server/modules/mails/services"
	profiles_services "github.com/guicostaarantes/psi-server/modules/profiles/services"
	translations_services "github.com/guicostaarantes/psi-server/modules/translations/services"
	treatments_services "github.com/guicostaarantes/psi-server/modules/treatments/services"
	users_services "github.com/guicostaarantes/psi-server/modules/users/services"
	"github.com/guicostaarantes/psi-server/utils/file_storage"
	"github.com/guicostaarantes/psi-server/utils/hash"
	"github.com/guicostaarantes/psi-server/utils/identifier"
	"github.com/guicostaarantes/psi-server/utils/mail"
	"github.com/guicostaarantes/psi-server/utils/match"
	"github.com/guicostaarantes/psi-server/utils/orm"
	"github.com/guicostaarantes/psi-server/utils/serializing"
	"github.com/guicostaarantes/psi-server/utils/token"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver receives all utils and registers all services within the application
type Resolver struct {
	OrmUtil                                   orm.IOrmUtil
	FileStorageUtil                           file_storage.IFileStorageUtil
	HashUtil                                  hash.IHashUtil
	IdentifierUtil                            identifier.IIdentifierUtil
	MailUtil                                  mail.IMailUtil
	MatchUtil                                 match.IMatchUtil
	SerializingUtil                           serializing.ISerializingUtil
	TokenUtil                                 token.ITokenUtil
	MaxAffinityNumber                         int64
	ScheduleIntervalDuration                  time.Duration
	ExpireAuthTokenDuration                   time.Duration
	ExpireResetTokenDuration                  time.Duration
	InterruptTreatmentCooldownDuration        time.Duration
	TopAffinitiesCooldownDuration             time.Duration
	askResetPasswordService                   *users_services.AskResetPasswordService
	assignTreatmentService                    *treatments_services.AssignTreatmentService
	authenticateUserService                   *users_services.AuthenticateUserService
	cancelAppointmentByPatientService         *appointments_services.CancelAppointmentByPatientService
	cancelAppointmentByPsychologistService    *appointments_services.CancelAppointmentByPsychologistService
	checkTreatmentCollisionService            *treatments_services.CheckTreatmentCollisionService
	confirmAppointmentByPatientService        *appointments_services.ConfirmAppointmentByPatientService
	confirmAppointmentByPsychologistService   *appointments_services.ConfirmAppointmentByPsychologistService
	createPendingAppointmentsService          *appointments_services.CreatePendingAppointmentsService
	createTreatmentService                    *treatments_services.CreateTreatmentService
	createUserService                         *users_services.CreateUserService
	createUserWithPasswordService             *users_services.CreateUserWithPasswordService
	deleteTreatmentService                    *treatments_services.DeleteTreatmentService
	editAppointmentByPatientService           *appointments_services.EditAppointmentByPatientService
	editAppointmentByPsychologistService      *appointments_services.EditAppointmentByPsychologistService
	finalizeTreatmentService                  *treatments_services.FinalizeTreatmentService
	getAgreementsByProfileIdService           *agreements_services.GetAgreementsByProfileIdService
	getAppointmentsOfPatientService           *appointments_services.GetAppointmentsOfPatientService
	getAppointmentsOfPsychologistService      *appointments_services.GetAppointmentsOfPsychologistService
	getCharacteristicsByIDService             *characteristics_services.GetCharacteristicsByIDService
	getCharacteristicsService                 *characteristics_services.GetCharacteristicsService
	getCooldownService                        *cooldowns_services.GetCooldownService
	getPatientByUserIDService                 *profiles_services.GetPatientByUserIDService
	getPatientService                         *profiles_services.GetPatientService
	getPatientTreatmentsService               *treatments_services.GetPatientTreatmentsService
	getPreferencesByIDService                 *characteristics_services.GetPreferencesByIDService
	getPsychologistByUserIDService            *profiles_services.GetPsychologistByUserIDService
	getPsychologistService                    *profiles_services.GetPsychologistService
	getPsychologistPendingTreatmentsService   *treatments_services.GetPsychologistPendingTreatmentsService
	getPsychologistPriceRangeOfferingsService *treatments_services.GetPsychologistPriceRangeOfferingsService
	getPsychologistTreatmentsService          *treatments_services.GetPsychologistTreatmentsService
	getTermsByProfileTypeService              *agreements_services.GetTermsByProfileTypeService
	getTopAffinitiesForPatientService         *characteristics_services.GetTopAffinitiesForPatientService
	getTranslationsService                    *translations_services.GetTranslationsService
	getTreatmentForPatientService             *treatments_services.GetTreatmentForPatientService
	getTreatmentForPsychologistService        *treatments_services.GetTreatmentForPsychologistService
	getTreatmentPriceRangeByNameService       *treatments_services.GetTreatmentPriceRangeByNameService
	getTreatmentPriceRangesService            *treatments_services.GetTreatmentPriceRangesService
	getUserByIDService                        *users_services.GetUserByIDService
	getUsersByRoleService                     *users_services.GetUsersByRoleService
	interruptTreatmentByPatientService        *treatments_services.InterruptTreatmentByPatientService
	interruptTreatmentByPsychologistService   *treatments_services.InterruptTreatmentByPsychologistService
	processPendingMailsService                *mails_services.ProcessPendingMailsService
	resetPasswordService                      *users_services.ResetPasswordService
	readFileService                           *files_services.ReadFileService
	saveCooldownService                       *cooldowns_services.SaveCooldownService
	setCharacteristicChoicesService           *characteristics_services.SetCharacteristicChoicesService
	setCharacteristicsService                 *characteristics_services.SetCharacteristicsService
	setPreferencesService                     *characteristics_services.SetPreferencesService
	setTopAffinitiesForPatientService         *characteristics_services.SetTopAffinitiesForPatientService
	setTranslationsService                    *translations_services.SetTranslationsService
	setTreatmentPriceRangesService            *treatments_services.SetTreatmentPriceRangesService
	updateTreatmentService                    *treatments_services.UpdateTreatmentService
	updateUserService                         *users_services.UpdateUserService
	uploadAvatarFileService                   *files_services.UploadAvatarFileService
	upsertAgreementService                    *agreements_services.UpsertAgreementService
	upsertPatientService                      *profiles_services.UpsertPatientService
	upsertPsychologistService                 *profiles_services.UpsertPsychologistService
	upsertTermService                         *agreements_services.UpsertTermService
	validateUserTokenService                  *users_services.ValidateUserTokenService
}

// AskResetPasswordService gets or sets the service with same name
func (r *Resolver) AskResetPasswordService() *users_services.AskResetPasswordService {
	if r.askResetPasswordService == nil {
		r.askResetPasswordService = &users_services.AskResetPasswordService{
			IdentifierUtil:           r.IdentifierUtil,
			OrmUtil:                  r.OrmUtil,
			TokenUtil:                r.TokenUtil,
			ExpireResetTokenDuration: r.ExpireResetTokenDuration,
		}
	}
	return r.askResetPasswordService
}

// AssignTreatmentService gets or sets the service with same name
func (r *Resolver) AssignTreatmentService() *treatments_services.AssignTreatmentService {
	if r.assignTreatmentService == nil {
		r.assignTreatmentService = &treatments_services.AssignTreatmentService{
			OrmUtil:            r.OrmUtil,
			GetCooldownService: r.GetCooldownService(),
		}
	}
	return r.assignTreatmentService
}

// AuthenticateUserService gets or sets the service with same name
func (r *Resolver) AuthenticateUserService() *users_services.AuthenticateUserService {
	if r.authenticateUserService == nil {
		r.authenticateUserService = &users_services.AuthenticateUserService{
			HashUtil:                r.HashUtil,
			OrmUtil:                 r.OrmUtil,
			SerializingUtil:         r.SerializingUtil,
			TokenUtil:               r.TokenUtil,
			ExpireAuthTokenDuration: r.ExpireAuthTokenDuration,
		}
	}
	return r.authenticateUserService
}

// CancelAppointmentByPatientService gets or sets the service with same name
func (r *Resolver) CancelAppointmentByPatientService() *appointments_services.CancelAppointmentByPatientService {
	if r.cancelAppointmentByPatientService == nil {
		r.cancelAppointmentByPatientService = &appointments_services.CancelAppointmentByPatientService{
			OrmUtil: r.OrmUtil,
		}
	}
	return r.cancelAppointmentByPatientService
}

// CancelAppointmentByPsychologistService gets or sets the service with same name
func (r *Resolver) CancelAppointmentByPsychologistService() *appointments_services.CancelAppointmentByPsychologistService {
	if r.cancelAppointmentByPsychologistService == nil {
		r.cancelAppointmentByPsychologistService = &appointments_services.CancelAppointmentByPsychologistService{
			OrmUtil: r.OrmUtil,
		}
	}
	return r.cancelAppointmentByPsychologistService
}

// CheckTreatmentCollisionService gets or sets the service with same name
func (r *Resolver) CheckTreatmentCollisionService() *treatments_services.CheckTreatmentCollisionService {
	if r.checkTreatmentCollisionService == nil {
		r.checkTreatmentCollisionService = &treatments_services.CheckTreatmentCollisionService{
			OrmUtil:                  r.OrmUtil,
			ScheduleIntervalDuration: r.ScheduleIntervalDuration,
		}
	}
	return r.checkTreatmentCollisionService
}

// ConfirmAppointmentByPatientService gets or sets the service with same name
func (r *Resolver) ConfirmAppointmentByPatientService() *appointments_services.ConfirmAppointmentByPatientService {
	if r.confirmAppointmentByPatientService == nil {
		r.confirmAppointmentByPatientService = &appointments_services.ConfirmAppointmentByPatientService{
			OrmUtil: r.OrmUtil,
		}
	}
	return r.confirmAppointmentByPatientService
}

// ConfirmAppointmentByPsychologistService gets or sets the service with same name
func (r *Resolver) ConfirmAppointmentByPsychologistService() *appointments_services.ConfirmAppointmentByPsychologistService {
	if r.confirmAppointmentByPsychologistService == nil {
		r.confirmAppointmentByPsychologistService = &appointments_services.ConfirmAppointmentByPsychologistService{
			OrmUtil: r.OrmUtil,
		}
	}
	return r.confirmAppointmentByPsychologistService
}

// CreatePendingAppointmentsService gets or sets the service with same name
func (r *Resolver) CreatePendingAppointmentsService() *appointments_services.CreatePendingAppointmentsService {
	if r.createPendingAppointmentsService == nil {
		r.createPendingAppointmentsService = &appointments_services.CreatePendingAppointmentsService{
			IdentifierUtil:           r.IdentifierUtil,
			OrmUtil:                  r.OrmUtil,
			ScheduleIntervalDuration: r.ScheduleIntervalDuration,
		}
	}
	return r.createPendingAppointmentsService
}

// CreateTreatmentService gets or sets the service with same name
func (r *Resolver) CreateTreatmentService() *treatments_services.CreateTreatmentService {
	if r.createTreatmentService == nil {
		r.createTreatmentService = &treatments_services.CreateTreatmentService{
			IdentifierUtil:                 r.IdentifierUtil,
			OrmUtil:                        r.OrmUtil,
			CheckTreatmentCollisionService: r.CheckTreatmentCollisionService(),
		}
	}
	return r.createTreatmentService
}

// CreateUserService gets or sets the service with same name
func (r *Resolver) CreateUserService() *users_services.CreateUserService {
	if r.createUserService == nil {
		r.createUserService = &users_services.CreateUserService{
			IdentifierUtil:          r.IdentifierUtil,
			MatchUtil:               r.MatchUtil,
			OrmUtil:                 r.OrmUtil,
			SerializingUtil:         r.SerializingUtil,
			TokenUtil:               r.TokenUtil,
			ExpireAuthTokenDuration: r.ExpireResetTokenDuration,
		}
	}
	return r.createUserService
}

// CreateUserWithPasswordService gets or sets the service with same name
func (r *Resolver) CreateUserWithPasswordService() *users_services.CreateUserWithPasswordService {
	if r.createUserWithPasswordService == nil {
		r.createUserWithPasswordService = &users_services.CreateUserWithPasswordService{
			HashUtil:        r.HashUtil,
			IdentifierUtil:  r.IdentifierUtil,
			MatchUtil:       r.MatchUtil,
			OrmUtil:         r.OrmUtil,
			SerializingUtil: r.SerializingUtil,
		}
	}
	return r.createUserWithPasswordService
}

// DeleteTreatmentService gets or sets the service with same name
func (r *Resolver) DeleteTreatmentService() *treatments_services.DeleteTreatmentService {
	if r.deleteTreatmentService == nil {
		r.deleteTreatmentService = &treatments_services.DeleteTreatmentService{
			OrmUtil: r.OrmUtil,
		}
	}
	return r.deleteTreatmentService
}

// EditAppointmentByPatientService gets or sets the service with same name
func (r *Resolver) EditAppointmentByPatientService() *appointments_services.EditAppointmentByPatientService {
	if r.editAppointmentByPatientService == nil {
		r.editAppointmentByPatientService = &appointments_services.EditAppointmentByPatientService{
			IdentifierUtil: r.IdentifierUtil,
			OrmUtil:        r.OrmUtil,
		}
	}
	return r.editAppointmentByPatientService
}

// EditAppointmentByPsychologistService gets or sets the service with same name
func (r *Resolver) EditAppointmentByPsychologistService() *appointments_services.EditAppointmentByPsychologistService {
	if r.editAppointmentByPsychologistService == nil {
		r.editAppointmentByPsychologistService = &appointments_services.EditAppointmentByPsychologistService{
			IdentifierUtil: r.IdentifierUtil,
			OrmUtil:        r.OrmUtil,
		}
	}
	return r.editAppointmentByPsychologistService
}

// FinalizeTreatmentService gets or sets the service with same name
func (r *Resolver) FinalizeTreatmentService() *treatments_services.FinalizeTreatmentService {
	if r.finalizeTreatmentService == nil {
		r.finalizeTreatmentService = &treatments_services.FinalizeTreatmentService{
			IdentifierUtil: r.IdentifierUtil,
			OrmUtil:        r.OrmUtil,
		}
	}
	return r.finalizeTreatmentService
}

// GetAgreementsByProfileIdService gets or sets the service with same name
func (r *Resolver) GetAgreementsByProfileIdService() *agreements_services.GetAgreementsByProfileIdService {
	if r.getAgreementsByProfileIdService == nil {
		r.getAgreementsByProfileIdService = &agreements_services.GetAgreementsByProfileIdService{
			OrmUtil: r.OrmUtil,
		}
	}
	return r.getAgreementsByProfileIdService
}

// GetAppointmentsOfPatientService gets or sets the service with same name
func (r *Resolver) GetAppointmentsOfPatientService() *appointments_services.GetAppointmentsOfPatientService {
	if r.getAppointmentsOfPatientService == nil {
		r.getAppointmentsOfPatientService = &appointments_services.GetAppointmentsOfPatientService{
			OrmUtil: r.OrmUtil,
		}
	}
	return r.getAppointmentsOfPatientService
}

// GetAppointmentsOfPsychologistService gets or sets the service with same name
func (r *Resolver) GetAppointmentsOfPsychologistService() *appointments_services.GetAppointmentsOfPsychologistService {
	if r.getAppointmentsOfPsychologistService == nil {
		r.getAppointmentsOfPsychologistService = &appointments_services.GetAppointmentsOfPsychologistService{
			OrmUtil: r.OrmUtil,
		}
	}
	return r.getAppointmentsOfPsychologistService
}

// GetCharacteristicsByIDService gets or sets the service with same name
func (r *Resolver) GetCharacteristicsByIDService() *characteristics_services.GetCharacteristicsByIDService {
	if r.getCharacteristicsByIDService == nil {
		r.getCharacteristicsByIDService = &characteristics_services.GetCharacteristicsByIDService{
			OrmUtil: r.OrmUtil,
		}
	}
	return r.getCharacteristicsByIDService
}

// GetCharacteristicsService gets or sets the service with same name
func (r *Resolver) GetCharacteristicsService() *characteristics_services.GetCharacteristicsService {
	if r.getCharacteristicsService == nil {
		r.getCharacteristicsService = &characteristics_services.GetCharacteristicsService{
			OrmUtil: r.OrmUtil,
		}
	}
	return r.getCharacteristicsService
}

// GetCooldownService gets or sets the service with same name
func (r *Resolver) GetCooldownService() *cooldowns_services.GetCooldownService {
	if r.getCooldownService == nil {
		r.getCooldownService = &cooldowns_services.GetCooldownService{
			OrmUtil: r.OrmUtil,
		}
	}
	return r.getCooldownService
}

// GetTranslationsService gets or sets the service with same name
func (r *Resolver) GetTranslationsService() *translations_services.GetTranslationsService {
	if r.getTranslationsService == nil {
		r.getTranslationsService = &translations_services.GetTranslationsService{
			OrmUtil: r.OrmUtil,
		}
	}
	return r.getTranslationsService
}

// GetTreatmentForPatientService gets or sets the service with same name
func (r *Resolver) GetTreatmentForPatientService() *treatments_services.GetTreatmentForPatientService {
	if r.getTreatmentForPatientService == nil {
		r.getTreatmentForPatientService = &treatments_services.GetTreatmentForPatientService{
			OrmUtil: r.OrmUtil,
		}
	}
	return r.getTreatmentForPatientService
}

// GetTreatmentForPsychologistService gets or sets the service with same name
func (r *Resolver) GetTreatmentForPsychologistService() *treatments_services.GetTreatmentForPsychologistService {
	if r.getTreatmentForPsychologistService == nil {
		r.getTreatmentForPsychologistService = &treatments_services.GetTreatmentForPsychologistService{
			OrmUtil: r.OrmUtil,
		}
	}
	return r.getTreatmentForPsychologistService
}

// GetPatientByUserIDService gets or sets the service with same name
func (r *Resolver) GetPatientByUserIDService() *profiles_services.GetPatientByUserIDService {
	if r.getPatientByUserIDService == nil {
		r.getPatientByUserIDService = &profiles_services.GetPatientByUserIDService{
			OrmUtil: r.OrmUtil,
		}
	}
	return r.getPatientByUserIDService
}

// GetPsychologistService gets or sets the service with same name
func (r *Resolver) GetPsychologistService() *profiles_services.GetPsychologistService {
	if r.getPsychologistService == nil {
		r.getPsychologistService = &profiles_services.GetPsychologistService{
			OrmUtil: r.OrmUtil,
		}
	}
	return r.getPsychologistService
}

// GetPatientTreatmentsService gets or sets the service with same name
func (r *Resolver) GetPatientTreatmentsService() *treatments_services.GetPatientTreatmentsService {
	if r.getPatientTreatmentsService == nil {
		r.getPatientTreatmentsService = &treatments_services.GetPatientTreatmentsService{
			OrmUtil: r.OrmUtil,
		}
	}
	return r.getPatientTreatmentsService
}

// GetPreferencesByIDService gets or sets the service with same name
func (r *Resolver) GetPreferencesByIDService() *characteristics_services.GetPreferencesByIDService {
	if r.getPreferencesByIDService == nil {
		r.getPreferencesByIDService = &characteristics_services.GetPreferencesByIDService{
			OrmUtil: r.OrmUtil,
		}
	}
	return r.getPreferencesByIDService
}

// GetPsychologistByUserIDService gets or sets the service with same name
func (r *Resolver) GetPsychologistByUserIDService() *profiles_services.GetPsychologistByUserIDService {
	if r.getPsychologistByUserIDService == nil {
		r.getPsychologistByUserIDService = &profiles_services.GetPsychologistByUserIDService{
			OrmUtil: r.OrmUtil,
		}
	}
	return r.getPsychologistByUserIDService
}

// GetPatientService gets or sets the service with same name
func (r *Resolver) GetPatientService() *profiles_services.GetPatientService {
	if r.getPatientService == nil {
		r.getPatientService = &profiles_services.GetPatientService{
			OrmUtil: r.OrmUtil,
		}
	}
	return r.getPatientService
}

// GetPsychologistPendingTreatmentsService gets or sets the service with same name
func (r *Resolver) GetPsychologistPendingTreatmentsService() *treatments_services.GetPsychologistPendingTreatmentsService {
	if r.getPsychologistPendingTreatmentsService == nil {
		r.getPsychologistPendingTreatmentsService = &treatments_services.GetPsychologistPendingTreatmentsService{
			OrmUtil: r.OrmUtil,
		}
	}
	return r.getPsychologistPendingTreatmentsService
}

// GetPsychologistPriceRangeOfferingsService gets or sets the service with same name
func (r *Resolver) GetPsychologistPriceRangeOfferingsService() *treatments_services.GetPsychologistPriceRangeOfferingsService {
	if r.getPsychologistPriceRangeOfferingsService == nil {
		r.getPsychologistPriceRangeOfferingsService = &treatments_services.GetPsychologistPriceRangeOfferingsService{
			OrmUtil: r.OrmUtil,
		}
	}
	return r.getPsychologistPriceRangeOfferingsService
}

// GetPsychologistTreatmentsService gets or sets the service with same name
func (r *Resolver) GetPsychologistTreatmentsService() *treatments_services.GetPsychologistTreatmentsService {
	if r.getPsychologistTreatmentsService == nil {
		r.getPsychologistTreatmentsService = &treatments_services.GetPsychologistTreatmentsService{
			OrmUtil: r.OrmUtil,
		}
	}
	return r.getPsychologistTreatmentsService
}

// GetTermsByProfileTypeService gets or sets the service with same name
func (r *Resolver) GetTermsByProfileTypeService() *agreements_services.GetTermsByProfileTypeService {
	if r.getTermsByProfileTypeService == nil {
		r.getTermsByProfileTypeService = &agreements_services.GetTermsByProfileTypeService{
			OrmUtil: r.OrmUtil,
		}
	}
	return r.getTermsByProfileTypeService
}

// GetTopAffinitiesForPatientService gets or sets the service with same name
func (r *Resolver) GetTopAffinitiesForPatientService() *characteristics_services.GetTopAffinitiesForPatientService {
	if r.getTopAffinitiesForPatientService == nil {
		r.getTopAffinitiesForPatientService = &characteristics_services.GetTopAffinitiesForPatientService{
			OrmUtil:                           r.OrmUtil,
			GetCooldownService:                r.GetCooldownService(),
			SetTopAffinitiesForPatientService: r.SetTopAffinitiesForPatientService(),
		}
	}
	return r.getTopAffinitiesForPatientService
}

// GetTreatmentPriceRangeByNameService gets or sets the service with same name
func (r *Resolver) GetTreatmentPriceRangeByNameService() *treatments_services.GetTreatmentPriceRangeByNameService {
	if r.getTreatmentPriceRangeByNameService == nil {
		r.getTreatmentPriceRangeByNameService = &treatments_services.GetTreatmentPriceRangeByNameService{
			OrmUtil: r.OrmUtil,
		}
	}
	return r.getTreatmentPriceRangeByNameService
}

// GetTreatmentPriceRangesService gets or sets the service with same name
func (r *Resolver) GetTreatmentPriceRangesService() *treatments_services.GetTreatmentPriceRangesService {
	if r.getTreatmentPriceRangesService == nil {
		r.getTreatmentPriceRangesService = &treatments_services.GetTreatmentPriceRangesService{
			OrmUtil: r.OrmUtil,
		}
	}
	return r.getTreatmentPriceRangesService
}

// GetUsersByRoleService gets or sets the service with same name
func (r *Resolver) GetUsersByRoleService() *users_services.GetUsersByRoleService {
	if r.getUsersByRoleService == nil {
		r.getUsersByRoleService = &users_services.GetUsersByRoleService{
			OrmUtil: r.OrmUtil,
		}
	}
	return r.getUsersByRoleService
}

// GetUserByIDService gets or sets the service with same name
func (r *Resolver) GetUserByIDService() *users_services.GetUserByIDService {
	if r.getUserByIDService == nil {
		r.getUserByIDService = &users_services.GetUserByIDService{
			OrmUtil: r.OrmUtil,
		}
	}
	return r.getUserByIDService
}

// InterruptTreatmentByPatientService gets or sets the service with same name
func (r *Resolver) InterruptTreatmentByPatientService() *treatments_services.InterruptTreatmentByPatientService {
	if r.interruptTreatmentByPatientService == nil {
		r.interruptTreatmentByPatientService = &treatments_services.InterruptTreatmentByPatientService{
			IdentifierUtil:      r.IdentifierUtil,
			OrmUtil:             r.OrmUtil,
			SaveCooldownService: r.SaveCooldownService(),
		}
	}
	return r.interruptTreatmentByPatientService
}

// InterruptTreatmentByPsychologistService gets or sets the service with same name
func (r *Resolver) InterruptTreatmentByPsychologistService() *treatments_services.InterruptTreatmentByPsychologistService {
	if r.interruptTreatmentByPsychologistService == nil {
		r.interruptTreatmentByPsychologistService = &treatments_services.InterruptTreatmentByPsychologistService{
			IdentifierUtil: r.IdentifierUtil,
			OrmUtil:        r.OrmUtil,
		}
	}
	return r.interruptTreatmentByPsychologistService
}

// ProcessPendingMailsService gets or sets the service with same name
func (r *Resolver) ProcessPendingMailsService() *mails_services.ProcessPendingMailsService {
	if r.processPendingMailsService == nil {
		r.processPendingMailsService = &mails_services.ProcessPendingMailsService{
			MailUtil: r.MailUtil,
			OrmUtil:  r.OrmUtil,
		}
	}
	return r.processPendingMailsService
}

// ReadFileService gets or sets the service with same name
func (r *Resolver) ReadFileService() *files_services.ReadFileService {
	if r.readFileService == nil {
		r.readFileService = &files_services.ReadFileService{
			FileStorageUtil: r.FileStorageUtil,
		}
	}
	return r.readFileService
}

// ResetPasswordService gets or sets the service with same name
func (r *Resolver) ResetPasswordService() *users_services.ResetPasswordService {
	if r.resetPasswordService == nil {
		r.resetPasswordService = &users_services.ResetPasswordService{
			HashUtil:  r.HashUtil,
			MatchUtil: r.MatchUtil,
			OrmUtil:   r.OrmUtil,
		}
	}
	return r.resetPasswordService
}

// SaveCooldownService gets or sets the service with same name
func (r *Resolver) SaveCooldownService() *cooldowns_services.SaveCooldownService {
	if r.saveCooldownService == nil {
		r.saveCooldownService = &cooldowns_services.SaveCooldownService{
			IdentifierUtil:                     r.IdentifierUtil,
			OrmUtil:                            r.OrmUtil,
			InterruptTreatmentCooldownDuration: r.InterruptTreatmentCooldownDuration,
			TopAffinitiesCooldownDuration:      r.TopAffinitiesCooldownDuration,
		}
	}
	return r.saveCooldownService
}

// SetCharacteristicChoicesService gets or sets the service with same name
func (r *Resolver) SetCharacteristicChoicesService() *characteristics_services.SetCharacteristicChoicesService {
	if r.setCharacteristicChoicesService == nil {
		r.setCharacteristicChoicesService = &characteristics_services.SetCharacteristicChoicesService{
			IdentifierUtil: r.IdentifierUtil,
			OrmUtil:        r.OrmUtil,
		}
	}
	return r.setCharacteristicChoicesService
}

// SetCharacteristicsService gets or sets the service with same name
func (r *Resolver) SetCharacteristicsService() *characteristics_services.SetCharacteristicsService {
	if r.setCharacteristicsService == nil {
		r.setCharacteristicsService = &characteristics_services.SetCharacteristicsService{
			IdentifierUtil: r.IdentifierUtil,
			OrmUtil:        r.OrmUtil,
		}
	}
	return r.setCharacteristicsService
}

// SetTranslationsService gets or sets the service with same name
func (r *Resolver) SetTranslationsService() *translations_services.SetTranslationsService {
	if r.setTranslationsService == nil {
		r.setTranslationsService = &translations_services.SetTranslationsService{
			OrmUtil: r.OrmUtil,
		}
	}
	return r.setTranslationsService
}

// SetPreferencesService gets or sets the service with same name
func (r *Resolver) SetPreferencesService() *characteristics_services.SetPreferencesService {
	if r.setPreferencesService == nil {
		r.setPreferencesService = &characteristics_services.SetPreferencesService{
			OrmUtil: r.OrmUtil,
		}
	}
	return r.setPreferencesService
}

// SetTopAffinitiesForPatientService gets or sets the service with same name
func (r *Resolver) SetTopAffinitiesForPatientService() *characteristics_services.SetTopAffinitiesForPatientService {
	if r.setTopAffinitiesForPatientService == nil {
		r.setTopAffinitiesForPatientService = &characteristics_services.SetTopAffinitiesForPatientService{
			IdentifierUtil:      r.IdentifierUtil,
			OrmUtil:             r.OrmUtil,
			MaxAffinityNumber:   r.MaxAffinityNumber,
			SaveCooldownService: r.SaveCooldownService(),
		}
	}
	return r.setTopAffinitiesForPatientService
}

// SetTreatmentPriceRangesService gets or sets the service with same name
func (r *Resolver) SetTreatmentPriceRangesService() *treatments_services.SetTreatmentPriceRangesService {
	if r.setTreatmentPriceRangesService == nil {
		r.setTreatmentPriceRangesService = &treatments_services.SetTreatmentPriceRangesService{
			IdentifierUtil: r.IdentifierUtil,
			OrmUtil:        r.OrmUtil,
		}
	}
	return r.setTreatmentPriceRangesService
}

// UpdateTreatmentService gets or sets the service with same name
func (r *Resolver) UpdateTreatmentService() *treatments_services.UpdateTreatmentService {
	if r.updateTreatmentService == nil {
		r.updateTreatmentService = &treatments_services.UpdateTreatmentService{
			IdentifierUtil:                 r.IdentifierUtil,
			OrmUtil:                        r.OrmUtil,
			CheckTreatmentCollisionService: r.CheckTreatmentCollisionService(),
		}
	}
	return r.updateTreatmentService
}

// UpdateUserService gets or sets the service with same name
func (r *Resolver) UpdateUserService() *users_services.UpdateUserService {
	if r.updateUserService == nil {
		r.updateUserService = &users_services.UpdateUserService{
			OrmUtil: r.OrmUtil,
		}
	}
	return r.updateUserService
}

// UploadAvatarFileService gets or sets the service with same name
func (r *Resolver) UploadAvatarFileService() *files_services.UploadAvatarFileService {
	if r.uploadAvatarFileService == nil {
		r.uploadAvatarFileService = &files_services.UploadAvatarFileService{
			FileStorageUtil: r.FileStorageUtil,
		}
	}
	return r.uploadAvatarFileService
}

// UpsertAgreementService gets or sets the service with same name
func (r *Resolver) UpsertAgreementService() *agreements_services.UpsertAgreementService {
	if r.upsertAgreementService == nil {
		r.upsertAgreementService = &agreements_services.UpsertAgreementService{
			IdentifierUtil: r.IdentifierUtil,
			OrmUtil:        r.OrmUtil,
		}
	}
	return r.upsertAgreementService
}

// UpsertPatientService gets or sets the service with same name
func (r *Resolver) UpsertPatientService() *profiles_services.UpsertPatientService {
	if r.upsertPatientService == nil {
		r.upsertPatientService = &profiles_services.UpsertPatientService{
			IdentifierUtil:          r.IdentifierUtil,
			OrmUtil:                 r.OrmUtil,
			UploadAvatarFileService: r.UploadAvatarFileService(),
		}
	}
	return r.upsertPatientService
}

// UpsertPsychologistService gets or sets the service with same name
func (r *Resolver) UpsertPsychologistService() *profiles_services.UpsertPsychologistService {
	if r.upsertPsychologistService == nil {
		r.upsertPsychologistService = &profiles_services.UpsertPsychologistService{
			IdentifierUtil:          r.IdentifierUtil,
			OrmUtil:                 r.OrmUtil,
			UploadAvatarFileService: r.UploadAvatarFileService(),
		}
	}
	return r.upsertPsychologistService
}

// UpsertTermService gets or sets the service with same name
func (r *Resolver) UpsertTermService() *agreements_services.UpsertTermService {
	if r.upsertTermService == nil {
		r.upsertTermService = &agreements_services.UpsertTermService{
			IdentifierUtil: r.IdentifierUtil,
			OrmUtil:        r.OrmUtil,
		}
	}
	return r.upsertTermService
}

// ValidateUserTokenService gets or sets the service with same name
func (r *Resolver) ValidateUserTokenService() *users_services.ValidateUserTokenService {
	if r.validateUserTokenService == nil {
		r.validateUserTokenService = &users_services.ValidateUserTokenService{
			OrmUtil:                 r.OrmUtil,
			SerializingUtil:         r.SerializingUtil,
			ExpireAuthTokenDuration: r.ExpireAuthTokenDuration,
		}
	}
	return r.validateUserTokenService
}
