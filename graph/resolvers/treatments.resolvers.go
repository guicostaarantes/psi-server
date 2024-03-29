package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/guicostaarantes/psi-server/graph/generated"
	profiles_models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	treatments_models "github.com/guicostaarantes/psi-server/modules/treatments/models"
)

func (r *mutationResolver) AssignTreatment(ctx context.Context, id string, priceRangeName string) (*bool, error) {
	userID := ctx.Value("userID").(string)

	servicePatient, servicePatientErr := r.GetPatientByUserIDService().Execute(userID)
	if servicePatientErr != nil {
		return nil, servicePatientErr
	}

	serviceErr := r.AssignTreatmentService().Execute(id, priceRangeName, servicePatient.ID)

	return nil, serviceErr
}

func (r *mutationResolver) CreateTreatment(ctx context.Context, input treatments_models.CreateTreatmentInput) (*bool, error) {
	userID := ctx.Value("userID").(string)

	servicePsy, servicePsyErr := r.GetPsychologistByUserIDService().Execute(userID)
	if servicePsyErr != nil {
		return nil, servicePsyErr
	}

	serviceErr := r.CreateTreatmentService().Execute(servicePsy.ID, input)

	return nil, serviceErr
}

func (r *mutationResolver) DeleteTreatment(ctx context.Context, id string, priceRangeName string) (*bool, error) {
	userID := ctx.Value("userID").(string)

	servicePsy, servicePsyErr := r.GetPsychologistByUserIDService().Execute(userID)
	if servicePsyErr != nil {
		return nil, servicePsyErr
	}

	serviceErr := r.DeleteTreatmentService().Execute(id, servicePsy.ID, priceRangeName)

	return nil, serviceErr
}

func (r *mutationResolver) InterruptTreatmentByPatient(ctx context.Context, id string, reason string) (*bool, error) {
	userID := ctx.Value("userID").(string)

	servicePatient, servicePatientErr := r.GetPatientByUserIDService().Execute(userID)
	if servicePatientErr != nil {
		return nil, servicePatientErr
	}

	serviceErr := r.InterruptTreatmentByPatientService().Execute(id, servicePatient.ID, reason)

	return nil, serviceErr
}

func (r *mutationResolver) InterruptTreatmentByPsychologist(ctx context.Context, id string, reason string) (*bool, error) {
	userID := ctx.Value("userID").(string)

	servicePsy, servicePsyErr := r.GetPsychologistByUserIDService().Execute(userID)
	if servicePsyErr != nil {
		return nil, servicePsyErr
	}

	serviceErr := r.InterruptTreatmentByPsychologistService().Execute(id, servicePsy.ID, reason)

	return nil, serviceErr
}

func (r *mutationResolver) FinalizeTreatment(ctx context.Context, id string) (*bool, error) {
	userID := ctx.Value("userID").(string)

	servicePsy, servicePsyErr := r.GetPsychologistByUserIDService().Execute(userID)
	if servicePsyErr != nil {
		return nil, servicePsyErr
	}

	serviceErr := r.FinalizeTreatmentService().Execute(id, servicePsy.ID)

	return nil, serviceErr
}

func (r *mutationResolver) SetTreatmentPriceRanges(ctx context.Context, input []*treatments_models.TreatmentPriceRange) (*bool, error) {
	serviceErr := r.SetTreatmentPriceRangesService().Execute(input)

	return nil, serviceErr
}

func (r *mutationResolver) UpdateTreatment(ctx context.Context, id string, input treatments_models.UpdateTreatmentInput) (*bool, error) {
	userID := ctx.Value("userID").(string)

	servicePsy, servicePsyErr := r.GetPsychologistByUserIDService().Execute(userID)
	if servicePsyErr != nil {
		return nil, servicePsyErr
	}

	serviceErr := r.UpdateTreatmentService().Execute(id, servicePsy.ID, input)

	return nil, serviceErr
}

func (r *patientTreatmentResolver) PriceRange(ctx context.Context, obj *treatments_models.GetPatientTreatmentsResponse) (*treatments_models.TreatmentPriceRange, error) {
	return r.GetTreatmentPriceRangeByNameService().Execute(obj.PriceRangeName)
}

func (r *patientTreatmentResolver) Psychologist(ctx context.Context, obj *treatments_models.GetPatientTreatmentsResponse) (*profiles_models.Psychologist, error) {
	return r.GetPsychologistService().Execute(obj.PsychologistID)
}

func (r *psychologistTreatmentResolver) PriceRange(ctx context.Context, obj *treatments_models.GetPsychologistTreatmentsResponse) (*treatments_models.TreatmentPriceRange, error) {
	return r.GetTreatmentPriceRangeByNameService().Execute(obj.PriceRangeName)
}

func (r *psychologistTreatmentResolver) Patient(ctx context.Context, obj *treatments_models.GetPsychologistTreatmentsResponse) (*profiles_models.Patient, error) {
	return r.GetPatientService().Execute(obj.PatientID)
}

func (r *queryResolver) TreatmentPriceRanges(ctx context.Context) ([]*treatments_models.TreatmentPriceRange, error) {
	return r.GetTreatmentPriceRangesService().Execute()
}

func (r *treatmentPriceRangeOfferingResolver) PriceRange(ctx context.Context, obj *treatments_models.TreatmentPriceRangeOffering) (*treatments_models.TreatmentPriceRange, error) {
	return r.GetTreatmentPriceRangeByNameService().Execute(obj.PriceRangeName)
}

// PatientTreatment returns generated.PatientTreatmentResolver implementation.
func (r *Resolver) PatientTreatment() generated.PatientTreatmentResolver {
	return &patientTreatmentResolver{r}
}

// PsychologistTreatment returns generated.PsychologistTreatmentResolver implementation.
func (r *Resolver) PsychologistTreatment() generated.PsychologistTreatmentResolver {
	return &psychologistTreatmentResolver{r}
}

// TreatmentPriceRangeOffering returns generated.TreatmentPriceRangeOfferingResolver implementation.
func (r *Resolver) TreatmentPriceRangeOffering() generated.TreatmentPriceRangeOfferingResolver {
	return &treatmentPriceRangeOfferingResolver{r}
}

type patientTreatmentResolver struct{ *Resolver }
type psychologistTreatmentResolver struct{ *Resolver }
type treatmentPriceRangeOfferingResolver struct{ *Resolver }
