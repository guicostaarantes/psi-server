package services

import (
	"strings"

	"github.com/guicostaarantes/psi-server/modules/characteristics/models"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// SetCharacteristicsService is a service that sets all possible patient characteristics
type SetCharacteristicsService struct {
	OrmUtil orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s SetCharacteristicsService) Execute(target models.CharacteristicTarget, input []*models.SetCharacteristicInput) error {

	newCharacteristics := []*models.Characteristic{}

	for _, char := range input {
		characteristic := models.Characteristic{
			Name:           char.Name,
			Type:           char.Type,
			Target:         target,
			PossibleValues: strings.Join(char.PossibleValues, ","),
		}

		newCharacteristics = append(newCharacteristics, &characteristic)
	}

	result := s.OrmUtil.Db().Delete(&models.Characteristic{}, "target = ?", target)
	if result.Error != nil {
		return result.Error
	}

	result = s.OrmUtil.Db().Create(&newCharacteristics)
	if result.Error != nil {
		return result.Error
	}

	return nil

}
