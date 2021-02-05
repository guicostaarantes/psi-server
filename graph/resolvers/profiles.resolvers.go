package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/guicostaarantes/psi-server/graph/generated"
	profiles_models "github.com/guicostaarantes/psi-server/modules/profiles/models"
)

func (r *mutationResolver) CreateOwnPsychologistProfile(ctx context.Context, input profiles_models.CreatePsychologistInput) (*bool, error) {
	userID := ctx.Value("userID").(string)

	serviceInput := &profiles_models.CreatePsychologistInput{
		UserID: userID,
	}

	mergeErr := r.MergeUtil.Merge(serviceInput, &input)
	if mergeErr != nil {
		return nil, mergeErr
	}

	serviceErr := r.CreatePsychologistService().Execute(serviceInput)
	if serviceErr != nil {
		return nil, serviceErr
	}

	return nil, nil
}

func (r *mutationResolver) CreatePsyCharacteristic(ctx context.Context, input profiles_models.CreatePsyCharacteristicInput) (*bool, error) {
	serviceErr := r.CreatePsyCharacteristicService().Execute(&input)
	if serviceErr != nil {
		return nil, serviceErr
	}

	return nil, nil
}

func (r *mutationResolver) SetOwnPsyCharacteristicChoice(ctx context.Context, input profiles_models.SetPsyCharacteristicChoiceInput) (*bool, error) {
	userID := ctx.Value("userID").(string)

	servicePsy, servicePsyErr := r.GetPsychologistByUserIDService().Execute(userID)
	if servicePsyErr != nil {
		return nil, servicePsyErr
	}

	serviceInput := &profiles_models.SetPsyCharacteristicChoiceInput{
		PsychologistID: servicePsy.ID,
	}

	mergeErr := r.MergeUtil.Merge(serviceInput, &input)
	if mergeErr != nil {
		return nil, mergeErr
	}

	serviceErr := r.SetPsyCharacteristicChoiceService().Execute(serviceInput)
	if serviceErr != nil {
		return nil, serviceErr
	}

	return nil, nil
}

func (r *mutationResolver) UpdateOwnPsychologistProfile(ctx context.Context, input profiles_models.UpdatePsychologistInput) (*bool, error) {
	userID := ctx.Value("userID").(string)

	servicePsy, servicePsyErr := r.GetPsychologistByUserIDService().Execute(userID)
	if servicePsyErr != nil {
		return nil, servicePsyErr
	}

	serviceErr := r.UpdatePsychologistService().Execute(servicePsy.ID, &input)
	if serviceErr != nil {
		return nil, serviceErr
	}

	return nil, nil
}

func (r *mutationResolver) UpdatePsyCharacteristic(ctx context.Context, id string, input profiles_models.UpdatePsyCharacteristicInput) (*bool, error) {
	serviceErr := r.UpdatePsyCharacteristicService().Execute(id, &input)
	if serviceErr != nil {
		return nil, serviceErr
	}

	return nil, nil
}

func (r *psychologistProfileResolver) Characteristics(ctx context.Context, obj *profiles_models.Psychologist) ([]*profiles_models.PsyCharacteristicChoiceResponse, error) {
	return r.GetPsyCharacteristicsByPsyIDService().Execute(obj.ID)
}

func (r *queryResolver) GetOwnPsychologistProfile(ctx context.Context) (*profiles_models.Psychologist, error) {
	userID := ctx.Value("userID").(string)
	return r.GetPsychologistByUserIDService().Execute(userID)
}

func (r *queryResolver) GetPsyCharacteristics(ctx context.Context) ([]*profiles_models.PsyCharacteristicResponse, error) {
	return r.GetPsyCharacteristicsService().Execute()
}

// PsychologistProfile returns generated.PsychologistProfileResolver implementation.
func (r *Resolver) PsychologistProfile() generated.PsychologistProfileResolver {
	return &psychologistProfileResolver{r}
}

type psychologistProfileResolver struct{ *Resolver }
