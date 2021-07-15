package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/guicostaarantes/psi-server/graph/generated"
	appointments_models "github.com/guicostaarantes/psi-server/modules/appointments/models"
	characteristics_models "github.com/guicostaarantes/psi-server/modules/characteristics/models"
	profiles_models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/modules/treatments/models"
)

func (r *mutationResolver) SetMyPatientCharacteristicChoices(ctx context.Context, input []*characteristics_models.SetCharacteristicChoiceInput) (*bool, error) {
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

func (r *mutationResolver) SetMyPatientPreferences(ctx context.Context, input []*characteristics_models.SetPreferenceInput) (*bool, error) {
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

func (r *mutationResolver) SetMyPsychologistCharacteristicChoices(ctx context.Context, input []*characteristics_models.SetCharacteristicChoiceInput) (*bool, error) {
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

func (r *mutationResolver) SetMyPsychologistPreferences(ctx context.Context, input []*characteristics_models.SetPreferenceInput) (*bool, error) {
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

func (r *mutationResolver) UpsertMyPatientProfile(ctx context.Context, input profiles_models.UpsertPatientInput) (*bool, error) {
	userID := ctx.Value("userID").(string)

	serviceInput := &profiles_models.UpsertPatientInput{
		UserID:    userID,
		FullName:  input.FullName,
		LikeName:  input.LikeName,
		BirthDate: input.BirthDate,
		City:      input.City,
	}

	serviceErr := r.UpsertPatientService().Execute(serviceInput)
	if serviceErr != nil {
		return nil, serviceErr
	}

	return nil, nil
}

func (r *mutationResolver) UpsertMyPsychologistProfile(ctx context.Context, input profiles_models.UpsertPsychologistInput) (*bool, error) {
	userID := ctx.Value("userID").(string)

	serviceInput := &profiles_models.UpsertPsychologistInput{
		UserID:    userID,
		FullName:  input.FullName,
		LikeName:  input.LikeName,
		BirthDate: input.BirthDate,
		City:      input.City,
	}

	serviceErr := r.UpsertPsychologistService().Execute(serviceInput)
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

func (r *patientProfileResolver) Treatments(ctx context.Context, obj *profiles_models.Patient) ([]*models.GetPatientTreatmentsResponse, error) {
	return r.GetPatientTreatmentsService().Execute(obj.ID)
}

func (r *patientProfileResolver) Appointments(ctx context.Context, obj *profiles_models.Patient) ([]*appointments_models.Appointment, error) {
	return r.GetAppointmentsOfPatientService().Execute(obj.ID)
}

func (r *psychologistProfileResolver) Characteristics(ctx context.Context, obj *profiles_models.Psychologist) ([]*characteristics_models.CharacteristicChoiceResponse, error) {
	return r.GetCharacteristicsByIDService().Execute(obj.ID)
}

func (r *psychologistProfileResolver) Preferences(ctx context.Context, obj *profiles_models.Psychologist) ([]*characteristics_models.PreferenceResponse, error) {
	return r.GetPreferencesByIDService().Execute(obj.ID)
}

func (r *psychologistProfileResolver) Treatments(ctx context.Context, obj *profiles_models.Psychologist) ([]*models.GetPsychologistTreatmentsResponse, error) {
	return r.GetPsychologistTreatmentsService().Execute(obj.ID)
}

func (r *psychologistProfileResolver) Appointments(ctx context.Context, obj *profiles_models.Psychologist) ([]*appointments_models.Appointment, error) {
	return r.GetAppointmentsOfPsychologistService().Execute(obj.ID)
}

func (r *publicPatientProfileResolver) Characteristics(ctx context.Context, obj *profiles_models.Patient) ([]*characteristics_models.CharacteristicChoiceResponse, error) {
	return r.GetCharacteristicsByIDService().Execute(obj.ID)
}

func (r *publicPsychologistProfileResolver) Characteristics(ctx context.Context, obj *profiles_models.Psychologist) ([]*characteristics_models.CharacteristicChoiceResponse, error) {
	return r.GetCharacteristicsByIDService().Execute(obj.ID)
}

func (r *publicPsychologistProfileResolver) PendingTreatments(ctx context.Context, obj *profiles_models.Psychologist) ([]*models.GetPsychologistTreatmentsResponse, error) {
	return r.GetPsychologistPendingTreatmentsService().Execute(obj.ID)
}

func (r *queryResolver) MyPatientProfile(ctx context.Context) (*profiles_models.Patient, error) {
	userID := ctx.Value("userID").(string)
	return r.GetPatientByUserIDService().Execute(userID)
}

func (r *queryResolver) MyPsychologistProfile(ctx context.Context) (*profiles_models.Psychologist, error) {
	userID := ctx.Value("userID").(string)
	return r.GetPsychologistByUserIDService().Execute(userID)
}

func (r *queryResolver) PatientProfile(ctx context.Context, id string) (*profiles_models.Patient, error) {
	return r.GetPatientService().Execute(id)
}

func (r *queryResolver) PsychologistProfile(ctx context.Context, id string) (*profiles_models.Psychologist, error) {
	return r.GetPsychologistService().Execute(id)
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
