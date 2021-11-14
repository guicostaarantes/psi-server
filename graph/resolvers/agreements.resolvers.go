package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	agreements_models "github.com/guicostaarantes/psi-server/modules/agreements/models"
)

func (r *mutationResolver) UpsertPatientAgreement(ctx context.Context, input agreements_models.UpsertAgreementInput) (*bool, error) {
	userID := ctx.Value("userID").(string)

	servicePatient, servicePatientErr := r.GetPatientByUserIDService().Execute(userID)
	if servicePatientErr != nil {
		return nil, servicePatientErr
	}

	serviceErr := r.UpsertAgreementService().Execute(servicePatient.ID, &input)

	return nil, serviceErr
}

func (r *mutationResolver) UpsertPsychologistAgreement(ctx context.Context, input agreements_models.UpsertAgreementInput) (*bool, error) {
	userID := ctx.Value("userID").(string)

	servicePsy, servicePsyErr := r.GetPsychologistByUserIDService().Execute(userID)
	if servicePsyErr != nil {
		return nil, servicePsyErr
	}

	serviceErr := r.UpsertAgreementService().Execute(servicePsy.ID, &input)

	return nil, serviceErr
}

func (r *mutationResolver) UpsertTerm(ctx context.Context, input agreements_models.Term) (*bool, error) {
	serviceErr := r.UpsertTermService().Execute(input.Name, input.Version, input.ProfileType, input.Active)

	return nil, serviceErr
}

func (r *queryResolver) PatientTerms(ctx context.Context) ([]*agreements_models.Term, error) {
	return r.GetTermsByProfileTypeService().Execute(agreements_models.Patient)
}

func (r *queryResolver) PsychologistTerms(ctx context.Context) ([]*agreements_models.Term, error) {
	return r.GetTermsByProfileTypeService().Execute(agreements_models.Psychologist)
}
