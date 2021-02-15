package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

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

func (r *mutationResolver) CreateOwnTreatment(ctx context.Context, input models.CreateTreatmentInput) (*bool, error) {
	userID := ctx.Value("userID").(string)

	servicePsy, servicePsyErr := r.GetPsychologistByUserIDService().Execute(userID)
	if servicePsyErr != nil {
		return nil, servicePsyErr
	}

	serviceErr := r.CreateTreatmentService().Execute(servicePsy.ID, input)

	return nil, serviceErr
}

func (r *mutationResolver) DeleteOwnTreatment(ctx context.Context, id string) (*bool, error) {
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

func (r *mutationResolver) FinalizeOwnTreatment(ctx context.Context, id string) (*bool, error) {
	userID := ctx.Value("userID").(string)

	servicePsy, servicePsyErr := r.GetPsychologistByUserIDService().Execute(userID)
	if servicePsyErr != nil {
		return nil, servicePsyErr
	}

	serviceErr := r.FinalizeTreatmentService().Execute(id, servicePsy.ID)

	return nil, serviceErr
}

func (r *mutationResolver) UpdateOwnTreatment(ctx context.Context, id string, input models.UpdateTreatmentInput) (*bool, error) {
	userID := ctx.Value("userID").(string)

	servicePsy, servicePsyErr := r.GetPsychologistByUserIDService().Execute(userID)
	if servicePsyErr != nil {
		return nil, servicePsyErr
	}

	serviceErr := r.UpdateTreatmentService().Execute(id, servicePsy.ID, input)

	return nil, serviceErr
}
