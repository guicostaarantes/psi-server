package services

import (
	translations_models "github.com/guicostaarantes/psi-server/modules/translations/models"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// GetTranslationsService is a service that gets translated translations
type GetTranslationsService struct {
	OrmUtil orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s GetTranslationsService) Execute(lang string, keys []string) ([]*translations_models.Translation, error) {

	translations := []*translations_models.Translation{}

	result := s.OrmUtil.Db().Where("lang = ? AND key IN ?", lang, keys).Order("key ASC").Find(&translations)
	if result.Error != nil {
		return nil, result.Error
	}

	return translations, nil

}
