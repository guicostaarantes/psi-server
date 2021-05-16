package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/guicostaarantes/psi-server/modules/translations/models"
)

func (r *mutationResolver) SetTranslations(ctx context.Context, lang string, input []*models.TranslationInput) (*bool, error) {
	serviceErr := r.SetTranslationsService().Execute(lang, input)
	if serviceErr != nil {
		return nil, serviceErr
	}

	return nil, nil
}

func (r *queryResolver) Translations(ctx context.Context, lang string, keys []string) ([]*models.Translation, error) {
	return r.GetTranslationsService().Execute(lang, keys)
}
