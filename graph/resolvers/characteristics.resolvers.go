package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/guicostaarantes/psi-server/graph/generated"
	characteristics_models "github.com/guicostaarantes/psi-server/modules/characteristics/models"
	profiles_models "github.com/guicostaarantes/psi-server/modules/profiles/models"
)

func (r *affinityResolver) Psychologist(ctx context.Context, obj *characteristics_models.Affinity) (*profiles_models.Psychologist, error) {
	return r.GetPsychologistService().Execute(obj.PsychologistID)
}

func (r *mutationResolver) SetTopAffinitiesForOwnPatient(ctx context.Context) (*bool, error) {
	userID := ctx.Value("userID").(string)

	servicePatient, servicePatientErr := r.GetPatientByUserIDService().Execute(userID)
	if servicePatientErr != nil {
		return nil, servicePatientErr
	}

	serviceErr := r.SetTopAffinitiesForPatientService().Execute(servicePatient.ID)

	return nil, serviceErr
}

func (r *queryResolver) GetPatientCharacteristics(ctx context.Context) ([]*characteristics_models.CharacteristicResponse, error) {
	return r.GetCharacteristicsService().Execute(characteristics_models.PatientTarget)
}

func (r *queryResolver) GetPsychologistCharacteristics(ctx context.Context) ([]*characteristics_models.CharacteristicResponse, error) {
	return r.GetCharacteristicsService().Execute(characteristics_models.PsychologistTarget)
}

func (r *queryResolver) GetTopAffinitiesForOwnPatient(ctx context.Context) ([]*characteristics_models.Affinity, error) {
	userID := ctx.Value("userID").(string)

	servicePatient, servicePatientErr := r.GetPatientByUserIDService().Execute(userID)
	if servicePatientErr != nil {
		return nil, servicePatientErr
	}

	return r.GetTopAffinitiesForPatientService().Execute(servicePatient.ID)
}

// Affinity returns generated.AffinityResolver implementation.
func (r *Resolver) Affinity() generated.AffinityResolver { return &affinityResolver{r} }

type affinityResolver struct{ *Resolver }
