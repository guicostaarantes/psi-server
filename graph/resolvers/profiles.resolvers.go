package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/guicostaarantes/psi-server/graph/generated"
	profiles_models "github.com/guicostaarantes/psi-server/modules/profiles/models"
)

func (r *mutationResolver) CreateOwnPatientProfile(ctx context.Context, input profiles_models.CreatePatientInput) (*bool, error) {
	userID := ctx.Value("userID").(string)

	serviceInput := &profiles_models.CreatePatientInput{
		UserID:    userID,
		BirthDate: input.BirthDate,
		City:      input.City,
	}

	serviceErr := r.CreatePatientService().Execute(serviceInput)
	if serviceErr != nil {
		return nil, serviceErr
	}

	return nil, nil
}

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

func (r *mutationResolver) SetOwnPatientCharacteristicChoices(ctx context.Context, input []*profiles_models.SetPatientCharacteristicChoiceInput) (*bool, error) {
	userID := ctx.Value("userID").(string)

	servicePsy, servicePsyErr := r.GetPatientByUserIDService().Execute(userID)
	if servicePsyErr != nil {
		return nil, servicePsyErr
	}

	serviceErr := r.SetPatientCharacteristicChoicesService().Execute(servicePsy.ID, input)
	if serviceErr != nil {
		return nil, serviceErr
	}

	return nil, nil
}

func (r *mutationResolver) SetOwnPsychologistCharacteristicChoices(ctx context.Context, input []*profiles_models.SetPsychologistCharacteristicChoiceInput) (*bool, error) {
	userID := ctx.Value("userID").(string)

	servicePsy, servicePsyErr := r.GetPsychologistByUserIDService().Execute(userID)
	if servicePsyErr != nil {
		return nil, servicePsyErr
	}

	serviceErr := r.SetPsychologistCharacteristicChoicesService().Execute(servicePsy.ID, input)
	if serviceErr != nil {
		return nil, serviceErr
	}

	return nil, nil
}

func (r *mutationResolver) SetPatientCharacteristics(ctx context.Context, input []*profiles_models.SetPatientCharacteristicInput) (*bool, error) {
	serviceErr := r.SetPatientCharacteristicsService().Execute(input)
	if serviceErr != nil {
		return nil, serviceErr
	}

	return nil, nil
}

func (r *mutationResolver) SetPsychologistCharacteristics(ctx context.Context, input []*profiles_models.SetPsychologistCharacteristicInput) (*bool, error) {
	serviceErr := r.SetPsychologistCharacteristicsService().Execute(input)
	if serviceErr != nil {
		return nil, serviceErr
	}

	return nil, nil
}

func (r *mutationResolver) UpdateOwnPatientProfile(ctx context.Context, input profiles_models.UpdatePatientInput) (*bool, error) {
	userID := ctx.Value("userID").(string)

	servicePsy, servicePsyErr := r.GetPatientByUserIDService().Execute(userID)
	if servicePsyErr != nil {
		return nil, servicePsyErr
	}

	serviceErr := r.UpdatePatientService().Execute(servicePsy.ID, &input)
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

func (r *patientProfileResolver) Characteristics(ctx context.Context, obj *profiles_models.Patient) ([]*profiles_models.PatientCharacteristicChoiceResponse, error) {
	return r.GetPatientCharacteristicsByPatientIDService().Execute(obj.ID)
}

func (r *psychologistProfileResolver) Characteristics(ctx context.Context, obj *profiles_models.Psychologist) ([]*profiles_models.PsychologistCharacteristicChoiceResponse, error) {
	return r.GetPsychologistCharacteristicsByPsyIDService().Execute(obj.ID)
}

func (r *queryResolver) GetOwnPatientProfile(ctx context.Context) (*profiles_models.Patient, error) {
	userID := ctx.Value("userID").(string)
	return r.GetPatientByUserIDService().Execute(userID)
}

func (r *queryResolver) GetOwnPsychologistProfile(ctx context.Context) (*profiles_models.Psychologist, error) {
	userID := ctx.Value("userID").(string)
	return r.GetPsychologistByUserIDService().Execute(userID)
}

func (r *queryResolver) GetPatientCharacteristics(ctx context.Context) ([]*profiles_models.PatientCharacteristicResponse, error) {
	return r.GetPatientCharacteristicsService().Execute()
}

func (r *queryResolver) GetPsychologistCharacteristics(ctx context.Context) ([]*profiles_models.PsychologistCharacteristicResponse, error) {
	return r.GetPsychologistCharacteristicsService().Execute()
}

// PatientProfile returns generated.PatientProfileResolver implementation.
func (r *Resolver) PatientProfile() generated.PatientProfileResolver {
	return &patientProfileResolver{r}
}

// PsychologistProfile returns generated.PsychologistProfileResolver implementation.
func (r *Resolver) PsychologistProfile() generated.PsychologistProfileResolver {
	return &psychologistProfileResolver{r}
}

type patientProfileResolver struct{ *Resolver }
type psychologistProfileResolver struct{ *Resolver }
