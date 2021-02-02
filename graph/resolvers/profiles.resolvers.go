package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/guicostaarantes/psi-server/graph/generated/model"
	profiles_models "github.com/guicostaarantes/psi-server/modules/profiles/models"
)

func (r *mutationResolver) CreateOwnPsychologistProfile(ctx context.Context, input model.CreateOwnPsychologistProfileInput) (*bool, error) {
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

func (r *mutationResolver) CreatePsyCharacteristic(ctx context.Context, input model.CreatePsyCharacteristicInput) (*bool, error) {
	serviceInput := &profiles_models.CreatePsyCharacteristicInput{}

	mergeErr := r.MergeUtil.Merge(serviceInput, &input)
	if mergeErr != nil {
		return nil, mergeErr
	}

	serviceErr := r.CreatePsyCharacteristicService().Execute(serviceInput)
	if serviceErr != nil {
		return nil, serviceErr
	}

	return nil, nil
}

func (r *mutationResolver) UpdateOwnPsychologistProfile(ctx context.Context, input model.UpdateOwnPsychologistProfileInput) (*bool, error) {
	userID := ctx.Value("userID").(string)

	servicePsy, servicePsyErr := r.GetPsychologistByUserIDService().Execute(userID)
	if servicePsyErr != nil {
		return nil, servicePsyErr
	}

	serviceInput := &profiles_models.UpdatePsychologistInput{}

	mergeErr := r.MergeUtil.Merge(serviceInput, &input)
	if mergeErr != nil {
		return nil, mergeErr
	}

	serviceErr := r.UpdatePsychologistService().Execute(servicePsy.ID, serviceInput)
	if serviceErr != nil {
		return nil, serviceErr
	}

	return nil, nil
}

func (r *mutationResolver) UpdatePsyCharacteristic(ctx context.Context, id string, input model.UpdatePsyCharacteristicInput) (*bool, error) {
	serviceInput := &profiles_models.UpdatePsyCharacteristicInput{}

	mergeErr := r.MergeUtil.Merge(serviceInput, &input)
	if mergeErr != nil {
		return nil, mergeErr
	}

	serviceErr := r.UpdatePsyCharacteristicService().Execute(id, serviceInput)
	if serviceErr != nil {
		return nil, serviceErr
	}

	return nil, nil
}
