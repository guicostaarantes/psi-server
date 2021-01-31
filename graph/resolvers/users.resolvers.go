package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/guicostaarantes/psi-server/graph/generated"
	"github.com/guicostaarantes/psi-server/graph/generated/model"
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

func (r *mutationResolver) CreatePatient(ctx context.Context, input model.CreatePatientInput) (*bool, error) {
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

func (r *mutationResolver) CreateUserWithPassword(ctx context.Context, input model.CreateUserInput) (*bool, error) {
	serviceErr := r.CreateUserWithPasswordService().Execute(&users_models.CreateUserWithPasswordInput{
		Email:     input.Email,
		Password:  input.Password,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Role:      string(input.Role),
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

func (r *mutationResolver) ResetPassword(ctx context.Context, input model.ResetPasswordInput) (*bool, error) {
	serviceErr := r.ResetPasswordService().Execute(&users_models.ResetPasswordInput{
		Token:    input.Token,
		Password: input.Password,
	})

	if serviceErr != nil {
		return nil, serviceErr
	}

	return nil, nil
}

func (r *mutationResolver) UpdateOwnUser(ctx context.Context, input model.UpdateOwnUserInput) (*bool, error) {
	userID := ctx.Value("userID").(string)

	serviceInput := &users_models.UpdateUserInput{}

	mergeErr := r.MergeUtil.Merge(serviceInput, &input)
	if mergeErr != nil {
		return nil, mergeErr
	}

	serviceErr := r.UpdateUserService().Execute(userID, serviceInput)
	if serviceErr != nil {
		return nil, serviceErr
	}

	return nil, nil
}

func (r *mutationResolver) UpdateUser(ctx context.Context, id string, input model.UpdateUserInput) (*bool, error) {
	serviceInput := &users_models.UpdateUserInput{}

	mergeErr := r.MergeUtil.Merge(serviceInput, &input)
	if mergeErr != nil {
		return nil, mergeErr
	}

	serviceErr := r.UpdateUserService().Execute(id, serviceInput)
	if serviceErr != nil {
		return nil, serviceErr
	}

	return nil, nil
}

func (r *queryResolver) AuthenticateUser(ctx context.Context, input model.AuthenticateUserInput) (*model.Token, error) {
	auth, err := r.AuthenticateUserService().Execute(&users_models.AuthenticateUserInput{
		Email:     input.Email,
		Password:  input.Password,
		IPAddress: input.IPAddress,
	})

	if err != nil {
		return nil, err
	}

	return &model.Token{
		Token:     auth.Token,
		ExpiresAt: auth.ExpiresAt,
	}, nil
}

func (r *queryResolver) GetOwnUser(ctx context.Context) (*model.User, error) {
	userID := ctx.Value("userID").(string)

	user := model.User{}

	serviceUser, serviceErr := r.GetUserByIdService().Execute(userID)
	if serviceErr != nil {
		return nil, serviceErr
	}

	mergeErr := r.MergeUtil.Merge(&user, serviceUser)
	if mergeErr != nil {
		return nil, mergeErr
	}

	return &user, nil
}

func (r *queryResolver) GetUser(ctx context.Context, id string) (*model.User, error) {
	user := model.User{}

	serviceUser, serviceErr := r.GetUserByIdService().Execute(id)
	if serviceErr != nil {
		return nil, serviceErr
	}

	mergeErr := r.MergeUtil.Merge(&user, serviceUser)
	if mergeErr != nil {
		return nil, mergeErr
	}

	fmt.Printf("%#v \n", user)

	return &user, nil
}

func (r *queryResolver) ListUsersByRole(ctx context.Context, role model.Role) ([]*model.User, error) {
	users := []*model.User{}

	serviceUsers, serviceErr := r.GetUsersByRoleService().Execute(string(role))
	if serviceErr != nil {
		return nil, serviceErr
	}

	mergeErr := r.MergeUtil.Merge(&users, serviceUsers)
	if mergeErr != nil {
		return nil, mergeErr
	}

	return users, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
