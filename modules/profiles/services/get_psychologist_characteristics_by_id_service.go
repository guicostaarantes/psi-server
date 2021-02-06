package services

import (
	"context"
	"strings"

	models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// GetPsychologistCharacteristicsByPsyIDService is a service that gets the user by userId
type GetPsychologistCharacteristicsByPsyIDService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s GetPsychologistCharacteristicsByPsyIDService) Execute(id string) ([]*models.PsychologistCharacteristicChoiceResponse, error) {

	charMap := map[string]*models.PsychologistCharacteristicChoiceResponse{}

	charCursor, findErr := s.DatabaseUtil.FindMany("psi_db", "psychologist_characteristics", map[string]interface{}{})
	if findErr != nil {
		return nil, findErr
	}

	defer charCursor.Close(context.Background())

	for charCursor.Next(context.Background()) {
		characteristic := models.PsychologistCharacteristic{}

		decodeErr := charCursor.Decode(&characteristic)
		if decodeErr != nil {
			return nil, decodeErr
		}

		characteristicResponse := models.PsychologistCharacteristicChoiceResponse{
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
		characteristic := models.PsychologistCharacteristicChoice{}

		decodeErr := choiceCursor.Decode(&characteristic)
		if decodeErr != nil {
			return nil, decodeErr
		}

		char := charMap[characteristic.CharacteristicName]

		char.Values = append(char.Values, characteristic.Value)
	}

	characteristics := []*models.PsychologistCharacteristicChoiceResponse{}

	for _, char := range charMap {
		characteristics = append(characteristics, char)
	}

	return characteristics, nil

}
