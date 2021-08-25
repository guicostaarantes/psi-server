package services

import (
	"github.com/guicostaarantes/psi-server/modules/translations/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// SetTranslationsService is a service that sets translations
type SetTranslationsService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s SetTranslationsService) Execute(lang string, input []*models.TranslationInput) error {

	translations := []interface{}{}

	for _, msg := range input {
		deleteErr := s.DatabaseUtil.DeleteMany("translations", map[string]interface{}{"key": msg.Key})
		if deleteErr != nil {
			return deleteErr
		}

		newTranslation := models.Translation{
			Lang:  lang,
			Key:   msg.Key,
			Value: msg.Value,
		}

		translations = append(translations, newTranslation)
	}

	writeErr := s.DatabaseUtil.InsertMany("translations", translations)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
