package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/guicostaarantes/psi-server/graph/generated"
	users_models "github.com/guicostaarantes/psi-server/modules/users/models"
)

func (r *mutationResolver) ActivateUser(ctx context.Context, id string) (*bool, error) {
	serviceErr := r.ActivateUserService().Execute(id, true)

	if serviceErr != nil {
		return nil, serviceErr
	}

	return nil, nil
}

func (r *mutationResolver) AskResetPassword(ctx context.Context, email string) (*bool, error) {
	serviceErr := r.AskResetPasswordService().Execute(email)

	if serviceErr != nil {
		return nil, serviceErr
	}

	return nil, nil
}

func (r *mutationResolver) CreatePatientUser(ctx context.Context, input users_models.CreateUserInput) (*bool, error) {
	serviceErr := r.CreateUserService().Execute(&users_models.CreateUserInput{
		Email:     input.Email,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Role:      "PATIENT",
	})

	if serviceErr != nil {
		return nil, serviceErr
	}

	return nil, nil
}

func (r *mutationResolver) CreatePsychologistUser(ctx context.Context, input users_models.CreateUserInput) (*bool, error) {
	serviceErr := r.CreateUserService().Execute(&users_models.CreateUserInput{
		Email:     input.Email,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Role:      "PSYCHOLOGIST",
	})

	if serviceErr != nil {
		return nil, serviceErr
	}

	return nil, nil
}

func (r *mutationResolver) CreateUserWithPassword(ctx context.Context, input users_models.CreateUserWithPasswordInput) (*bool, error) {
	serviceErr := r.CreateUserWithPasswordService().Execute(&users_models.CreateUserWithPasswordInput{
		Email:     input.Email,
		Password:  input.Password,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Role:      input.Role,
	})

	if serviceErr != nil {
		return nil, serviceErr
	}

	return nil, nil
}

func (r *mutationResolver) DeactivateUser(ctx context.Context, id string) (*bool, error) {
	serviceErr := r.ActivateUserService().Execute(id, false)

	if serviceErr != nil {
		return nil, serviceErr
	}

	return nil, nil
}

func (r *mutationResolver) ResetPassword(ctx context.Context, input users_models.ResetPasswordInput) (*bool, error) {
	serviceErr := r.ResetPasswordService().Execute(&users_models.ResetPasswordInput{
		Token:    input.Token,
		Password: input.Password,
	})

	if serviceErr != nil {
		return nil, serviceErr
	}

	return nil, nil
}

func (r *mutationResolver) UpdateOwnUser(ctx context.Context, input users_models.UpdateUserInput) (*bool, error) {
	userID := ctx.Value("userID").(string)

	serviceErr := r.UpdateUserService().Execute(userID, &input)
	if serviceErr != nil {
		return nil, serviceErr
	}

	return nil, nil
}

func (r *mutationResolver) UpdateUser(ctx context.Context, id string, input users_models.UpdateUserInput) (*bool, error) {
	serviceErr := r.UpdateUserService().Execute(id, &input)

	if serviceErr != nil {
		return nil, serviceErr
	}

	return nil, nil
}

func (r *queryResolver) AuthenticateUser(ctx context.Context, input users_models.AuthenticateUserInput) (*users_models.Authentication, error) {
	return r.AuthenticateUserService().Execute(&users_models.AuthenticateUserInput{
		Email:     input.Email,
		Password:  input.Password,
		IPAddress: input.IPAddress,
	})
}

func (r *queryResolver) GetOwnUser(ctx context.Context) (*users_models.User, error) {
	userID := ctx.Value("userID").(string)
	return r.GetUserByIDService().Execute(userID)
}

func (r *queryResolver) GetUser(ctx context.Context, id string) (*users_models.User, error) {
	return r.GetUserByIDService().Execute(id)
}

func (r *queryResolver) ListUsersByRole(ctx context.Context, role users_models.Role) ([]*users_models.User, error) {
	return r.GetUsersByRoleService().Execute(string(role))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
