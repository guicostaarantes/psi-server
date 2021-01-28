package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/guicostaarantes/psi-server/graph/generated/model"
	users_models "github.com/guicostaarantes/psi-server/modules/users/models"
)

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
