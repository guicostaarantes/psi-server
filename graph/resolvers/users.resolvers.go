package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	generated1 "github.com/guicostaarantes/psi-server/graph/generated"
	users_models "github.com/guicostaarantes/psi-server/modules/users/models"
)

func (r *mutationResolver) AskResetPassword(ctx context.Context, email string) (*bool, error) {
	serviceErr := r.AskResetPasswordService().Execute(email)

	if serviceErr != nil {
		return nil, serviceErr
	}

	return nil, nil
}

func (r *mutationResolver) CreatePatientUser(ctx context.Context, input users_models.CreateUserInput) (*bool, error) {
	serviceErr := r.CreateUserService().Execute(&users_models.CreateUserInput{
		Email: input.Email,
		Role:  "PATIENT",
	})

	if serviceErr != nil {
		return nil, serviceErr
	}

	return nil, nil
}

func (r *mutationResolver) CreatePsychologistUser(ctx context.Context, input users_models.CreateUserInput) (*bool, error) {
	serviceErr := r.CreateUserService().Execute(&users_models.CreateUserInput{
		Email: input.Email,
		Role:  "PSYCHOLOGIST",
	})

	if serviceErr != nil {
		return nil, serviceErr
	}

	return nil, nil
}

func (r *mutationResolver) CreateUserWithPassword(ctx context.Context, input users_models.CreateUserWithPasswordInput) (*bool, error) {
	serviceErr := r.CreateUserWithPasswordService().Execute(&users_models.CreateUserWithPasswordInput{
		Email:    input.Email,
		Password: input.Password,
		Role:     input.Role,
	})

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

func (r *mutationResolver) UpdateUser(ctx context.Context, id string, input users_models.UpdateUserInput) (*bool, error) {
	serviceErr := r.UpdateUserService().Execute(id, &input)

	if serviceErr != nil {
		return nil, serviceErr
	}

	return nil, nil
}

func (r *queryResolver) AuthenticateUser(ctx context.Context, input users_models.AuthenticateUserInput) (*users_models.Authentication, error) {
	return r.AuthenticateUserService().Execute(&users_models.AuthenticateUserInput{
		Email:    input.Email,
		Password: input.Password,
	})
}

func (r *queryResolver) MyUser(ctx context.Context) (*users_models.User, error) {
	userID := ctx.Value("userID").(string)
	return r.GetUserByIDService().Execute(userID)
}

func (r *queryResolver) User(ctx context.Context, id string) (*users_models.User, error) {
	return r.GetUserByIDService().Execute(id)
}

func (r *queryResolver) UsersByRole(ctx context.Context, role users_models.Role) ([]*users_models.User, error) {
	return r.GetUsersByRoleService().Execute(string(role))
}

// Mutation returns generated1.MutationResolver implementation.
func (r *Resolver) Mutation() generated1.MutationResolver { return &mutationResolver{r} }

// Query returns generated1.QueryResolver implementation.
func (r *Resolver) Query() generated1.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
