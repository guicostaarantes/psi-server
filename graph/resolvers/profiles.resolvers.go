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
		UserID:    userID,
		BirthDate: input.BirthDate,
		City:      input.City,
	}

	serviceErr := r.CreatePsychologistService().Execute(serviceInput)
	if serviceErr != nil {
		return nil, serviceErr
	}

	return nil, nil
}

func (r *mutationResolver) CreatePsychologistCharacteristic(ctx context.Context, input profiles_models.CreatePsychologistCharacteristicInput) (*bool, error) {
	serviceErr := r.CreatePsychologistCharacteristicService().Execute(&input)
	if serviceErr != nil {
		return nil, serviceErr
	}

	return nil, nil
}

func (r *mutationResolver) SetOwnPsychologistCharacteristicChoice(ctx context.Context, input profiles_models.SetPsychologistCharacteristicChoiceInput) (*bool, error) {
	userID := ctx.Value("userID").(string)

	servicePsy, servicePsyErr := r.GetPsychologistByUserIDService().Execute(userID)
	if servicePsyErr != nil {
		return nil, servicePsyErr
	}

	serviceInput := &profiles_models.SetPsychologistCharacteristicChoiceInput{
		PsychologistID:     servicePsy.ID,
		CharacteristicName: input.CharacteristicName,
		Values:             input.Values,
	}

	serviceErr := r.SetPsychologistCharacteristicChoiceService().Execute(serviceInput)
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

func (r *mutationResolver) UpdatePsychologistCharacteristic(ctx context.Context, id string, input profiles_models.UpdatePsychologistCharacteristicInput) (*bool, error) {
	serviceErr := r.UpdatePsychologistCharacteristicService().Execute(id, &input)
	if serviceErr != nil {
		return nil, serviceErr
	}

	return nil, nil
}

func (r *psychologistProfileResolver) Characteristics(ctx context.Context, obj *profiles_models.Psychologist) ([]*profiles_models.PsychologistCharacteristicChoiceResponse, error) {
	return r.GetPsychologistCharacteristicsByPsyIDService().Execute(obj.ID)
}

func (r *queryResolver) GetOwnPsychologistProfile(ctx context.Context) (*profiles_models.Psychologist, error) {
	userID := ctx.Value("userID").(string)
	return r.GetPsychologistByUserIDService().Execute(userID)
}

func (r *queryResolver) GetPsychologistCharacteristics(ctx context.Context) ([]*profiles_models.PsychologistCharacteristicResponse, error) {
	return r.GetPsychologistCharacteristicsService().Execute()
}

// PsychologistProfile returns generated.PsychologistProfileResolver implementation.
func (r *Resolver) PsychologistProfile() generated.PsychologistProfileResolver {
	return &psychologistProfileResolver{r}
}

type psychologistProfileResolver struct{ *Resolver }
