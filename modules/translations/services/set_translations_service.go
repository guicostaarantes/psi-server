package services

import (
	translations_models "github.com/guicostaarantes/psi-server/modules/translations/models"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// SetTranslationsService is a service that sets translations
type SetTranslationsService struct {
	OrmUtil orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s SetTranslationsService) Execute(lang string, input []*translations_models.TranslationInput) error {

	keys := []string{}
	inputs := map[string]string{}

	for _, msg := range input {
		keys = append(keys, msg.Key)
		inputs[msg.Key] = msg.Value
	}

	translationResults := []*translations_models.Translation{}

	result := s.OrmUtil.Db().Where("lang = ? AND key IN ?", lang, keys).Find(&translationResults)
	if result.Error != nil {
		return result.Error
	}

	translations := map[string]*translations_models.Translation{}

	for _, trans := range translationResults {
		translations[trans.Key] = trans
	}

	for _, key := range keys {
		if _, exists := translations[key]; exists {
			translations[key].Value = inputs[key]
			result := s.OrmUtil.Db().Save(translations[key])
			if result.Error != nil {
				return result.Error
			}
		} else {
			newTranslation := translations_models.Translation{
				Lang:  lang,
				Key:   key,
				Value: inputs[key],
			}
			result := s.OrmUtil.Db().Create(&newTranslation)
			if result.Error != nil {
				return result.Error
			}
		}
	}

	return nil

}
