package services

import (
	"context"
	"strings"

	"github.com/guicostaarantes/psi-server/modules/characteristics/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// GetCharacteristicsService is a service that gets all possible characteristic based on the target
type GetCharacteristicsService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s GetCharacteristicsService) Execute(target models.CharacteristicTarget) ([]*models.CharacteristicResponse, error) {

	characteristics := []*models.CharacteristicResponse{}

	cursor, findErr := s.DatabaseUtil.FindMany("psi_db", "characteristics", map[string]interface{}{"target": string(target)})
	if findErr != nil {
		return nil, findErr
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		characteristic := models.Characteristic{}

		decodeErr := cursor.Decode(&characteristic)
		if decodeErr != nil {
			return nil, decodeErr
		}

		characteristicResponse := models.CharacteristicResponse{
			Name:           characteristic.Name,
			Type:           characteristic.Type,
			PossibleValues: strings.Split(characteristic.PossibleValues, ","),
		}

		characteristics = append(characteristics, &characteristicResponse)
	}

	return characteristics, nil

}
