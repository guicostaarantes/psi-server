package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/guicostaarantes/psi-server/graph/generated"
	models1 "github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/modules/schedule/models"
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

func (r *mutationResolver) ConfirmAppointment(ctx context.Context, id string) (*bool, error) {
	userID := ctx.Value("userID").(string)

	servicePsy, servicePsyErr := r.GetPsychologistByUserIDService().Execute(userID)
	if servicePsyErr != nil {
		return nil, servicePsyErr
	}

	serviceErr := r.ConfirmAppointmentService().Execute(id, servicePsy.ID)

	return nil, serviceErr
}

func (r *mutationResolver) DenyAppointment(ctx context.Context, id string, reason string) (*bool, error) {
	userID := ctx.Value("userID").(string)

	servicePsy, servicePsyErr := r.GetPsychologistByUserIDService().Execute(userID)
	if servicePsyErr != nil {
		return nil, servicePsyErr
	}

	serviceErr := r.DenyAppointmentService().Execute(id, servicePsy.ID, reason)

	return nil, serviceErr
}

func (r *mutationResolver) ProposeAppointment(ctx context.Context, input models.ProposeAppointmentInput) (*bool, error) {
	userID := ctx.Value("userID").(string)

	servicePatient, servicePatientErr := r.GetPatientByUserIDService().Execute(userID)
	if servicePatientErr != nil {
		return nil, servicePatientErr
	}

	serviceErr := r.ProposeAppointmentService().Execute(servicePatient.ID, input)

	return nil, serviceErr
}

func (r *patientAppointmentResolver) Psychologist(ctx context.Context, obj *models.Appointment) (*models1.Psychologist, error) {
	return r.GetPsychologistService().Execute(obj.PsychologistID)
}

func (r *psychologistAppointmentResolver) Patient(ctx context.Context, obj *models.Appointment) (*models1.Patient, error) {
	return r.GetPatientService().Execute(obj.PatientID)
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
