package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/guicostaarantes/psi-server/graph/generated"
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

func (r *mutationResolver) SetOwnPsyCharacteristicChoice(ctx context.Context, input model.SetOwnPsyCharacteristicChoiceInput) (*bool, error) {
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

func (r *psychologistProfileResolver) Characteristics(ctx context.Context, obj *model.PsychologistProfile) ([]*model.PsyCharacteristicChoice, error) {
	characteristics := []*model.PsyCharacteristicChoice{}

	serviceChars, serviceErr := r.GetPsyCharacteristicsByPsyIDService().Execute(obj.ID)
	if serviceErr != nil {
		return nil, serviceErr
	}

	mergeErr := r.MergeUtil.Merge(&characteristics, serviceChars)
	if mergeErr != nil {
		return nil, mergeErr
	}

	return characteristics, nil
}

func (r *queryResolver) GetOwnPsychologistProfile(ctx context.Context) (*model.PsychologistProfile, error) {
	userID := ctx.Value("userID").(string)

	profile := model.PsychologistProfile{}

	psy, serviceErr := r.GetPsychologistByUserIDService().Execute(userID)
	if serviceErr != nil {
		return nil, serviceErr
	}

	mergeErr := r.MergeUtil.Merge(&profile, psy)
	if mergeErr != nil {
		return nil, mergeErr
	}

	return &profile, nil
}

func (r *queryResolver) GetPsyCharacteristics(ctx context.Context) ([]*model.PsyCharacteristic, error) {
	characteristics, serviceErr := r.GetPsyCharacteristicsService().Execute()
	if serviceErr != nil {
		return nil, serviceErr
	}

	response := []*model.PsyCharacteristic{}

	for _, char := range characteristics {
		resp := &model.PsyCharacteristic{}

		mergeErr := r.MergeUtil.Merge(resp, char)
		if mergeErr != nil {
			return nil, mergeErr
		}

		response = append(response, resp)
	}

	return response, nil
}

// PsychologistProfile returns generated.PsychologistProfileResolver implementation.
func (r *Resolver) PsychologistProfile() generated.PsychologistProfileResolver {
	return &psychologistProfileResolver{r}
}

type psychologistProfileResolver struct{ *Resolver }
