package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/guicostaarantes/psi-server/graph/generated"
	appointments_models "github.com/guicostaarantes/psi-server/modules/appointments/models"
	treatments_models "github.com/guicostaarantes/psi-server/modules/treatments/models"
)

func (r *mutationResolver) CancelAppointmentByPatient(ctx context.Context, id string, reason string) (*bool, error) {
	userID := ctx.Value("userID").(string)

	servicePatient, servicePatientErr := r.GetPatientByUserIDService().Execute(userID)
	if servicePatientErr != nil {
		return nil, servicePatientErr
	}

	serviceErr := r.CancelAppointmentByPatientService().Execute(id, servicePatient.ID, reason)

	return nil, serviceErr
}

func (r *mutationResolver) CancelAppointmentByPsychologist(ctx context.Context, id string, reason string) (*bool, error) {
	userID := ctx.Value("userID").(string)

	servicePsy, servicePsyErr := r.GetPsychologistByUserIDService().Execute(userID)
	if servicePsyErr != nil {
		return nil, servicePsyErr
	}

	serviceErr := r.CancelAppointmentByPsychologistService().Execute(id, servicePsy.ID, reason)

	return nil, serviceErr
}

func (r *mutationResolver) ConfirmAppointmentByPatient(ctx context.Context, id string) (*bool, error) {
	userID := ctx.Value("userID").(string)

	servicePatient, servicePatientErr := r.GetPatientByUserIDService().Execute(userID)
	if servicePatientErr != nil {
		return nil, servicePatientErr
	}

	serviceErr := r.ConfirmAppointmentByPatientService().Execute(id, servicePatient.ID)

	return nil, serviceErr
}

func (r *mutationResolver) ConfirmAppointmentByPsychologist(ctx context.Context, id string) (*bool, error) {
	userID := ctx.Value("userID").(string)

	servicePsy, servicePsyErr := r.GetPsychologistByUserIDService().Execute(userID)
	if servicePsyErr != nil {
		return nil, servicePsyErr
	}

	serviceErr := r.ConfirmAppointmentByPsychologistService().Execute(id, servicePsy.ID)

	return nil, serviceErr
}

func (r *mutationResolver) CreatePendingAppointments(ctx context.Context) (*bool, error) {
	serviceErr := r.CreatePendingAppointmentsService().Execute()

	return nil, serviceErr
}

func (r *mutationResolver) EditAppointmentByPatient(ctx context.Context, id string, input appointments_models.EditAppointmentByPatientInput) (*bool, error) {
	userID := ctx.Value("userID").(string)

	servicePatient, servicePatientErr := r.GetPatientByUserIDService().Execute(userID)
	if servicePatientErr != nil {
		return nil, servicePatientErr
	}

	serviceErr := r.EditAppointmentByPatientService().Execute(id, servicePatient.ID, input)

	return nil, serviceErr
}

func (r *mutationResolver) EditAppointmentByPsychologist(ctx context.Context, id string, input appointments_models.EditAppointmentByPsychologistInput) (*bool, error) {
	userID := ctx.Value("userID").(string)

	servicePsy, servicePsyErr := r.GetPsychologistByUserIDService().Execute(userID)
	if servicePsyErr != nil {
		return nil, servicePsyErr
	}

	serviceErr := r.EditAppointmentByPsychologistService().Execute(id, servicePsy.ID, input)

	return nil, serviceErr
}

func (r *patientAppointmentResolver) PriceRange(ctx context.Context, obj *appointments_models.Appointment) (*treatments_models.TreatmentPriceRange, error) {
	return r.GetTreatmentPriceRangeByNameService().Execute(obj.PriceRangeName)
}

func (r *patientAppointmentResolver) Treatment(ctx context.Context, obj *appointments_models.Appointment) (*treatments_models.GetPatientTreatmentsResponse, error) {
	return r.GetTreatmentForPatientService().Execute(obj.TreatmentID)
}

func (r *psychologistAppointmentResolver) PriceRange(ctx context.Context, obj *appointments_models.Appointment) (*treatments_models.TreatmentPriceRange, error) {
	return r.GetTreatmentPriceRangeByNameService().Execute(obj.PriceRangeName)
}

func (r *psychologistAppointmentResolver) Treatment(ctx context.Context, obj *appointments_models.Appointment) (*treatments_models.GetPsychologistTreatmentsResponse, error) {
	return r.GetTreatmentForPsychologistService().Execute(obj.TreatmentID)
}

// PatientAppointment returns generated.PatientAppointmentResolver implementation.
func (r *Resolver) PatientAppointment() generated.PatientAppointmentResolver {
	return &patientAppointmentResolver{r}
}

// PsychologistAppointment returns generated.PsychologistAppointmentResolver implementation.
func (r *Resolver) PsychologistAppointment() generated.PsychologistAppointmentResolver {
	return &psychologistAppointmentResolver{r}
}

type patientAppointmentResolver struct{ *Resolver }
type psychologistAppointmentResolver struct{ *Resolver }
