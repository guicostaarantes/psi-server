package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/guicostaarantes/psi-server/graph/generated"
	characteristics_models "github.com/guicostaarantes/psi-server/modules/characteristics/models"
	profiles_models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	schedule_models "github.com/guicostaarantes/psi-server/modules/schedule/models"
)

func (r *mutationResolver) CreateOwnPatientProfile(ctx context.Context, input profiles_models.CreatePatientInput) (*bool, error) {
	userID := ctx.Value("userID").(string)

	serviceInput := &profiles_models.CreatePatientInput{
		UserID:    userID,
		FullName:  input.FullName,
		LikeName:  input.LikeName,
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
		FullName:  input.FullName,
		LikeName:  input.LikeName,
		BirthDate: input.BirthDate,
		City:      input.City,
	}

	serviceErr := r.CreatePsychologistService().Execute(serviceInput)
	if serviceErr != nil {
		return nil, serviceErr
	}

	return nil, nil
}

func (r *mutationResolver) SetOwnPatientCharacteristicChoices(ctx context.Context, input []*characteristics_models.SetCharacteristicChoiceInput) (*bool, error) {
	userID := ctx.Value("userID").(string)

	servicePatient, servicePatientErr := r.GetPatientByUserIDService().Execute(userID)
	if servicePatientErr != nil {
		return nil, servicePatientErr
	}

	serviceErr := r.SetCharacteristicChoicesService().Execute(servicePatient.ID, input)
	if serviceErr != nil {
		return nil, serviceErr
	}

	return nil, nil
}

func (r *mutationResolver) SetOwnPatientPreferences(ctx context.Context, input []*characteristics_models.SetPreferenceInput) (*bool, error) {
	userID := ctx.Value("userID").(string)

	servicePatient, servicePatientErr := r.GetPatientByUserIDService().Execute(userID)
	if servicePatientErr != nil {
		return nil, servicePatientErr
	}

	serviceErr := r.SetPreferencesService().Execute(servicePatient.ID, input)
	if serviceErr != nil {
		return nil, serviceErr
	}

	return nil, nil
}

func (r *mutationResolver) SetOwnPsychologistCharacteristicChoices(ctx context.Context, input []*characteristics_models.SetCharacteristicChoiceInput) (*bool, error) {
	userID := ctx.Value("userID").(string)

	servicePsy, servicePsyErr := r.GetPsychologistByUserIDService().Execute(userID)
	if servicePsyErr != nil {
		return nil, servicePsyErr
	}

	serviceErr := r.SetCharacteristicChoicesService().Execute(servicePsy.ID, input)
	if serviceErr != nil {
		return nil, serviceErr
	}

	return nil, nil
}

func (r *mutationResolver) SetOwnPsychologistPreferences(ctx context.Context, input []*characteristics_models.SetPreferenceInput) (*bool, error) {
	userID := ctx.Value("userID").(string)

	servicePsy, servicePsyErr := r.GetPsychologistByUserIDService().Execute(userID)
	if servicePsyErr != nil {
		return nil, servicePsyErr
	}

	serviceErr := r.SetPreferencesService().Execute(servicePsy.ID, input)
	if serviceErr != nil {
		return nil, serviceErr
	}

	return nil, nil
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

func (r *patientProfileResolver) Characteristics(ctx context.Context, obj *profiles_models.Patient) ([]*characteristics_models.CharacteristicChoiceResponse, error) {
	return r.GetCharacteristicsByIDService().Execute(obj.ID)
}

func (r *patientProfileResolver) Preferences(ctx context.Context, obj *profiles_models.Patient) ([]*characteristics_models.PreferenceResponse, error) {
	return r.GetPreferencesByIDService().Execute(obj.ID)
}

func (r *patientProfileResolver) Slots(ctx context.Context, obj *profiles_models.Patient) ([]*schedule_models.GetPatientSlotsResponse, error) {
	return r.GetPatientSlotsService().Execute(obj.ID)
}

func (r *psychologistProfileResolver) Characteristics(ctx context.Context, obj *profiles_models.Psychologist) ([]*characteristics_models.CharacteristicChoiceResponse, error) {
	return r.GetCharacteristicsByIDService().Execute(obj.ID)
}

func (r *psychologistProfileResolver) Preferences(ctx context.Context, obj *profiles_models.Psychologist) ([]*characteristics_models.PreferenceResponse, error) {
	return r.GetPreferencesByIDService().Execute(obj.ID)
}

func (r *psychologistProfileResolver) Slots(ctx context.Context, obj *profiles_models.Psychologist) ([]*schedule_models.GetPsychologistSlotsResponse, error) {
	return r.GetPsychologistSlotsService().Execute(obj.ID)
}

func (r *publicPatientProfileResolver) Characteristics(ctx context.Context, obj *profiles_models.Patient) ([]*characteristics_models.CharacteristicChoiceResponse, error) {
	return r.GetCharacteristicsByIDService().Execute(obj.ID)
}

func (r *publicPsychologistProfileResolver) Characteristics(ctx context.Context, obj *profiles_models.Psychologist) ([]*characteristics_models.CharacteristicChoiceResponse, error) {
	return r.GetCharacteristicsByIDService().Execute(obj.ID)
}

func (r *queryResolver) GetOwnPatientProfile(ctx context.Context) (*profiles_models.Patient, error) {
	userID := ctx.Value("userID").(string)
	return r.GetPatientByUserIDService().Execute(userID)
}

func (r *queryResolver) GetOwnPsychologistProfile(ctx context.Context) (*profiles_models.Psychologist, error) {
	userID := ctx.Value("userID").(string)
	return r.GetPsychologistByUserIDService().Execute(userID)
}

func (r *queryResolver) GetPatientCharacteristics(ctx context.Context) ([]*characteristics_models.CharacteristicResponse, error) {
	return r.GetCharacteristicsService().Execute(characteristics_models.PatientTarget)
}

func (r *queryResolver) GetPatientProfile(ctx context.Context, id string) (*profiles_models.Patient, error) {
	return r.GetPatientByUserIDService().Execute(id)
}

func (r *queryResolver) GetPsychologistCharacteristics(ctx context.Context) ([]*characteristics_models.CharacteristicResponse, error) {
	return r.GetCharacteristicsService().Execute(characteristics_models.PsychologistTarget)
}

func (r *queryResolver) GetPsychologistProfile(ctx context.Context, id string) (*profiles_models.Psychologist, error) {
	return r.GetPsychologistByUserIDService().Execute(id)
}

// PatientProfile returns generated.PatientProfileResolver implementation.
func (r *Resolver) PatientProfile() generated.PatientProfileResolver {
	return &patientProfileResolver{r}
}

// PsychologistProfile returns generated.PsychologistProfileResolver implementation.
func (r *Resolver) PsychologistProfile() generated.PsychologistProfileResolver {
	return &psychologistProfileResolver{r}
}

// PublicPatientProfile returns generated.PublicPatientProfileResolver implementation.
func (r *Resolver) PublicPatientProfile() generated.PublicPatientProfileResolver {
	return &publicPatientProfileResolver{r}
}

// PublicPsychologistProfile returns generated.PublicPsychologistProfileResolver implementation.
func (r *Resolver) PublicPsychologistProfile() generated.PublicPsychologistProfileResolver {
	return &publicPsychologistProfileResolver{r}
}

type patientProfileResolver struct{ *Resolver }
type psychologistProfileResolver struct{ *Resolver }
type publicPatientProfileResolver struct{ *Resolver }
type publicPsychologistProfileResolver struct{ *Resolver }
