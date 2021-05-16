package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/guicostaarantes/psi-server/graph/generated"
	models1 "github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/modules/treatments/models"
)

func (r *mutationResolver) AssignTreatment(ctx context.Context, id string) (*bool, error) {
	userID := ctx.Value("userID").(string)

	servicePatient, servicePatientErr := r.GetPatientByUserIDService().Execute(userID)
	if servicePatientErr != nil {
		return nil, servicePatientErr
	}

	serviceErr := r.AssignTreatmentService().Execute(id, servicePatient.ID)

	return nil, serviceErr
}

func (r *mutationResolver) CreateTreatment(ctx context.Context, input models.CreateTreatmentInput) (*bool, error) {
	userID := ctx.Value("userID").(string)

	servicePsy, servicePsyErr := r.GetPsychologistByUserIDService().Execute(userID)
	if servicePsyErr != nil {
		return nil, servicePsyErr
	}

	serviceErr := r.CreateTreatmentService().Execute(servicePsy.ID, input)

	return nil, serviceErr
}

func (r *mutationResolver) DeleteTreatment(ctx context.Context, id string) (*bool, error) {
	userID := ctx.Value("userID").(string)

	servicePsy, servicePsyErr := r.GetPsychologistByUserIDService().Execute(userID)
	if servicePsyErr != nil {
		return nil, servicePsyErr
	}

	serviceErr := r.DeleteTreatmentService().Execute(id, servicePsy.ID)

	return nil, serviceErr
}

func (r *mutationResolver) InterruptTreatmentByPatient(ctx context.Context, id string, reason string) (*bool, error) {
	userID := ctx.Value("userID").(string)

	servicePatient, servicePatientErr := r.GetPatientByUserIDService().Execute(userID)
	if servicePatientErr != nil {
		return nil, servicePatientErr
	}

	serviceErr := r.InterruptTreatmentByPatientService().Execute(id, servicePatient.ID, reason)

	return nil, serviceErr
}

func (r *mutationResolver) InterruptTreatmentByPsychologist(ctx context.Context, id string, reason string) (*bool, error) {
	userID := ctx.Value("userID").(string)

	servicePsy, servicePsyErr := r.GetPsychologistByUserIDService().Execute(userID)
	if servicePsyErr != nil {
		return nil, servicePsyErr
	}

	serviceErr := r.InterruptTreatmentByPsychologistService().Execute(id, servicePsy.ID, reason)

	return nil, serviceErr
}

func (r *mutationResolver) FinalizeTreatment(ctx context.Context, id string) (*bool, error) {
	userID := ctx.Value("userID").(string)

	servicePsy, servicePsyErr := r.GetPsychologistByUserIDService().Execute(userID)
	if servicePsyErr != nil {
		return nil, servicePsyErr
	}

	serviceErr := r.FinalizeTreatmentService().Execute(id, servicePsy.ID)

	return nil, serviceErr
}

func (r *mutationResolver) UpdateTreatment(ctx context.Context, id string, input models.UpdateTreatmentInput) (*bool, error) {
	userID := ctx.Value("userID").(string)

	servicePsy, servicePsyErr := r.GetPsychologistByUserIDService().Execute(userID)
	if servicePsyErr != nil {
		return nil, servicePsyErr
	}

	serviceErr := r.UpdateTreatmentService().Execute(id, servicePsy.ID, input)

	return nil, serviceErr
}

func (r *patientTreatmentResolver) Psychologist(ctx context.Context, obj *models.GetPatientTreatmentsResponse) (*models1.Psychologist, error) {
	return r.GetPsychologistService().Execute(obj.PsychologistID)
}

func (r *psychologistTreatmentResolver) Patient(ctx context.Context, obj *models.GetPsychologistTreatmentsResponse) (*models1.Patient, error) {
	return r.GetPatientService().Execute(obj.PatientID)
}

// PatientTreatment returns generated.PatientTreatmentResolver implementation.
func (r *Resolver) PatientTreatment() generated.PatientTreatmentResolver {
	return &patientTreatmentResolver{r}
}

// PsychologistTreatment returns generated.PsychologistTreatmentResolver implementation.
func (r *Resolver) PsychologistTreatment() generated.PsychologistTreatmentResolver {
	return &psychologistTreatmentResolver{r}
}

type patientTreatmentResolver struct{ *Resolver }
type psychologistTreatmentResolver struct{ *Resolver }
