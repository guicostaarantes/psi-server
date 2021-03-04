package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/guicostaarantes/psi-server/modules/messages/models"
)

func (r *mutationResolver) SetMessages(ctx context.Context, lang string, input []*models.MessageInput) (*bool, error) {
	serviceErr := r.SetMessagesService().Execute(lang, input)
	if serviceErr != nil {
		return nil, serviceErr
	}

	return nil, nil
}

func (r *queryResolver) GetMessages(ctx context.Context, lang string, keys []string) ([]*models.Message, error) {
	return r.GetMessagesService().Execute(lang, keys)
}
