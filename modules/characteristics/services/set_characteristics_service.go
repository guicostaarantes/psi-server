package characteristcs_services

import (
	"strings"

	characteristics_models "github.com/guicostaarantes/psi-server/modules/characteristics/models"
	"github.com/guicostaarantes/psi-server/utils/identifier"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// SetCharacteristicsService is a service that sets all possible patient characteristics
type SetCharacteristicsService struct {
	IdentifierUtil identifier.IIdentifierUtil
	OrmUtil        orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s SetCharacteristicsService) Execute(target characteristics_models.CharacteristicTarget, input []*characteristics_models.SetCharacteristicInput) error {

	currentCharacteristics := []*characteristics_models.Characteristic{}

	result := s.OrmUtil.Db().Where("target = ?", target).Find(&currentCharacteristics)
	if result.Error != nil {
		return result.Error
	}

	currentChars := map[string]*characteristics_models.Characteristic{}

	for _, char := range currentCharacteristics {
		currentChars[char.Name] = char
	}

	for _, char := range input {
		if _, exists := currentChars[char.Name]; exists {

			currentChars[char.Name].Type = char.Type
			currentChars[char.Name].PossibleValues = strings.Join(char.PossibleValues, ",")

			result := s.OrmUtil.Db().Save(currentChars[char.Name])
			if result.Error != nil {
				return result.Error
			}

			delete(currentChars, char.Name)

		} else {

			_, charID, charIDErr := s.IdentifierUtil.GenerateIdentifier()
			if charIDErr != nil {
				return charIDErr
			}

			result := s.OrmUtil.Db().Create(&characteristics_models.Characteristic{
				ID:             charID,
				Name:           char.Name,
				Type:           char.Type,
				Target:         target,
				PossibleValues: strings.Join(char.PossibleValues, ","),
			})
			if result.Error != nil {
				return result.Error
			}

		}

	}

	// Deleting remaining chars
	for _, char := range currentChars {
		result := s.OrmUtil.Db().Delete(&char)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil

}
