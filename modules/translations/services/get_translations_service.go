package services

import (
	"context"

	"github.com/guicostaarantes/psi-server/modules/translations/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// GetTranslationsService is a service that gets translated translations
type GetTranslationsService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s GetTranslationsService) Execute(lang string, keys []string) ([]*models.Translation, error) {

	translations := []*models.Translation{}

	cursor, findErr := s.DatabaseUtil.FindMany("translations", map[string]interface{}{"lang": lang})
	if findErr != nil {
		return nil, findErr
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		msg := models.Translation{}

		decodeErr := cursor.Decode(&msg)
		if decodeErr != nil {
			return nil, decodeErr
		}

		for _, key := range keys {
			if key == msg.Key {
				translations = append(translations, &msg)
			}
		}
	}

	return translations, nil

}
