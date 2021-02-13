package services

import (
	"context"
	"strings"

	models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// GetCharacteristicsByIDService is a service that gets the characteristics of a profile based on its id
type GetCharacteristicsByIDService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s GetCharacteristicsByIDService) Execute(id string) ([]*models.CharacteristicChoiceResponse, error) {

	characteristics := []*models.CharacteristicChoiceResponse{}

	charCursor, findErr := s.DatabaseUtil.FindMany("psi_db", "characteristics", map[string]interface{}{})
	if findErr != nil {
		return nil, findErr
	}

	defer charCursor.Close(context.Background())

	for charCursor.Next(context.Background()) {
		characteristic := models.Characteristic{}

		decodeErr := charCursor.Decode(&characteristic)
		if decodeErr != nil {
			return nil, decodeErr
		}

		characteristicResponse := models.CharacteristicChoiceResponse{
			Name:           characteristic.Name,
			Type:           characteristic.Type,
			SelectedValues: []string{},
			PossibleValues: strings.Split(characteristic.PossibleValues, ","),
		}

		characteristics = append(characteristics, &characteristicResponse)
	}

	choiceCursor, findErr := s.DatabaseUtil.FindMany("psi_db", "characteristic_choices", map[string]interface{}{"profileId": id})
	if findErr != nil {
		return nil, findErr
	}

	defer choiceCursor.Close(context.Background())

	for choiceCursor.Next(context.Background()) {
		choice := models.CharacteristicChoice{}

		decodeErr := choiceCursor.Decode(&choice)
		if decodeErr != nil {
			return nil, decodeErr
		}

		for _, char := range characteristics {
			if char.Name == choice.CharacteristicName {
				char.SelectedValues = append(char.SelectedValues, choice.SelectedValue)
			}
		}
	}

	return characteristics, nil

}
