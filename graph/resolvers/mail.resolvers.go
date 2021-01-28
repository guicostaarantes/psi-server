package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
)

func (r *mutationResolver) ProcessPendingMail(ctx context.Context) (*bool, error) {
	serviceErr := ProcessPendingMailsService.Execute()

	if serviceErr != nil {
		return nil, serviceErr
	}

	return nil, nil
}
