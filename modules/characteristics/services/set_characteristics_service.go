package services

import (
	"strings"

	"github.com/guicostaarantes/psi-server/modules/characteristics/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// SetCharacteristicsService is a service that sets all possible patient characteristics
type SetCharacteristicsService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s SetCharacteristicsService) Execute(target models.CharacteristicTarget, input []*models.SetCharacteristicInput) error {

	newCharacteristics := []interface{}{}

	for _, char := range input {
		characteristic := models.Characteristic{
			Name:           char.Name,
			Type:           char.Type,
			Target:         target,
			PossibleValues: strings.Join(char.PossibleValues, ","),
		}

		newCharacteristics = append(newCharacteristics, characteristic)
	}

	deleteErr := s.DatabaseUtil.DeleteMany("characteristics", map[string]interface{}{"target": string(target)})
	if deleteErr != nil {
		return deleteErr
	}

	writeErr := s.DatabaseUtil.InsertMany("characteristics", newCharacteristics)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
