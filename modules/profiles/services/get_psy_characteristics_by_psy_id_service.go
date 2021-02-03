package services

import (
	"context"
	"strings"

	models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// GetPsyCharacteristicsByPsyIDService is a service that gets the user by userId
type GetPsyCharacteristicsByPsyIDService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s GetPsyCharacteristicsByPsyIDService) Execute(id string) ([]*models.PsyCharacteristicChoiceResponse, error) {

	charMap := map[string]*models.PsyCharacteristicChoiceResponse{}

	charCursor, findErr := s.DatabaseUtil.FindMany("psi_db", "psychologist_characteristics", map[string]interface{}{})
	if findErr != nil {
		return nil, findErr
	}

	defer charCursor.Close(context.Background())

	for charCursor.Next(context.Background()) {
		characteristic := models.PsyCharacteristic{}

		decodeErr := charCursor.Decode(&characteristic)
		if decodeErr != nil {
			return nil, decodeErr
		}

		characteristicResponse := models.PsyCharacteristicChoiceResponse{
			ID:             characteristic.ID,
			Name:           characteristic.Name,
			Many:           characteristic.Many,
			Values:         []string{},
			PossibleValues: strings.Split(characteristic.PossibleValues, ","),
		}

		charMap[characteristic.Name] = &characteristicResponse
	}

	choiceCursor, findErr := s.DatabaseUtil.FindMany("psi_db", "psychologist_characteristic_choices", map[string]interface{}{"psychologistId": id})
	if findErr != nil {
		return nil, findErr
	}

	defer choiceCursor.Close(context.Background())

	for choiceCursor.Next(context.Background()) {
		characteristic := models.PsyCharacteristicChoice{}

		decodeErr := choiceCursor.Decode(&characteristic)
		if decodeErr != nil {
			return nil, decodeErr
		}

		char := charMap[characteristic.CharacteristicName]

		char.Values = append(char.Values, characteristic.Value)
	}

	characteristics := []*models.PsyCharacteristicChoiceResponse{}

	for _, char := range charMap {
		characteristics = append(characteristics, char)
	}

	return characteristics, nil

}
