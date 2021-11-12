package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/guicostaarantes/psi-server/graph/generated"
	characteristics_models "github.com/guicostaarantes/psi-server/modules/characteristics/models"
	profiles_models "github.com/guicostaarantes/psi-server/modules/profiles/models"
)

func (r *affinityResolver) CreatedAt(ctx context.Context, obj *characteristics_models.Affinity) (int64, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *affinityResolver) Psychologist(ctx context.Context, obj *characteristics_models.Affinity) (*profiles_models.Psychologist, error) {
	return r.GetPsychologistService().Execute(obj.PsychologistID)
}

func (r *mutationResolver) SetPatientCharacteristics(ctx context.Context, input []*characteristics_models.SetCharacteristicInput) (*bool, error) {
	serviceErr := r.SetCharacteristicsService().Execute(characteristics_models.PatientTarget, input)
	if serviceErr != nil {
		return nil, serviceErr
	}

	return nil, nil
}

func (r *mutationResolver) SetPsychologistCharacteristics(ctx context.Context, input []*characteristics_models.SetCharacteristicInput) (*bool, error) {
	serviceErr := r.SetCharacteristicsService().Execute(characteristics_models.PsychologistTarget, input)
	if serviceErr != nil {
		return nil, serviceErr
	}

	return nil, nil
}

func (r *queryResolver) PatientCharacteristics(ctx context.Context) ([]*characteristics_models.CharacteristicResponse, error) {
	return r.GetCharacteristicsService().Execute(characteristics_models.PatientTarget)
}

func (r *queryResolver) PsychologistCharacteristics(ctx context.Context) ([]*characteristics_models.CharacteristicResponse, error) {
	return r.GetCharacteristicsService().Execute(characteristics_models.PsychologistTarget)
}

func (r *queryResolver) MyPatientTopAffinities(ctx context.Context) ([]*characteristics_models.Affinity, error) {
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
