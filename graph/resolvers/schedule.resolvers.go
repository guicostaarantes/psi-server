package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/guicostaarantes/psi-server/modules/schedule/models"
)

func (r *mutationResolver) ConfirmAppointment(ctx context.Context, id string) (*bool, error) {
	userID := ctx.Value("userID").(string)

	servicePsy, servicePsyErr := r.GetPsychologistByUserIDService().Execute(userID)
	if servicePsyErr != nil {
		return nil, servicePsyErr
	}

	serviceErr := r.ConfirmAppointmentService().Execute(id, servicePsy.ID)

	return nil, serviceErr
}

func (r *mutationResolver) DenyAppointment(ctx context.Context, id string) (*bool, error) {
	userID := ctx.Value("userID").(string)

	servicePsy, servicePsyErr := r.GetPsychologistByUserIDService().Execute(userID)
	if servicePsyErr != nil {
		return nil, servicePsyErr
	}

	serviceErr := r.DenyAppointmentService().Execute(id, servicePsy.ID)

	return nil, serviceErr
}

func (r *mutationResolver) ProposeAppointment(ctx context.Context, input models.ProposeAppointmentInput) (*bool, error) {
	userID := ctx.Value("userID").(string)

	servicePatient, servicePatientErr := r.GetPatientByUserIDService().Execute(userID)
	if servicePatientErr != nil {
		return nil, servicePatientErr
	}

	serviceErr := r.ProposeAppointmentService().Execute(servicePatient.ID, input)

	return nil, serviceErr
}

func (r *mutationResolver) SetOwnAvailability(ctx context.Context, input []*models.SetAvailabilityInput) (*bool, error) {
	userID := ctx.Value("userID").(string)

	servicePsy, servicePsyErr := r.GetPsychologistByUserIDService().Execute(userID)
	if servicePsyErr != nil {
		return nil, servicePsyErr
	}

	serviceErr := r.SetAvailabilityService().Execute(servicePsy.ID, input)

	return nil, serviceErr
}

func (r *queryResolver) GetOwnAvailability(ctx context.Context) ([]*models.AvailabilityResponse, error) {
	userID := ctx.Value("userID").(string)

	servicePsy, servicePsyErr := r.GetPsychologistByUserIDService().Execute(userID)
	if servicePsyErr != nil {
		return nil, servicePsyErr
	}

	return r.GetAvailabilityService().Execute(servicePsy.ID)
}

func (r *queryResolver) GetPsychologistAvailability(ctx context.Context, id string) ([]*models.AvailabilityResponse, error) {
	servicePsy, servicePsyErr := r.GetPsychologistByUserIDService().Execute(id)
	if servicePsyErr != nil {
		return nil, servicePsyErr
	}

	return r.GetAvailabilityService().Execute(servicePsy.ID)
}
