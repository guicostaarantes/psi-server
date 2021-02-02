package profiles_services

import (
	"context"
	"strings"

	models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// GetPsyCharacteristicsService is a service that gets the user by userId
type GetPsyCharacteristicsService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s GetPsyCharacteristicsService) Execute() ([]*models.PsyCharacteristicResponse, error) {

	characteristics := []*models.PsyCharacteristicResponse{}

	cursor, findErr := s.DatabaseUtil.FindMany("psi_db", "psychologist_characteristics", map[string]interface{}{})
	if findErr != nil {
		return nil, findErr
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		characteristic := models.PsyCharacteristic{}

		decodeErr := cursor.Decode(&characteristic)
		if decodeErr != nil {
			return nil, decodeErr
		}

		characteristicResponse := models.PsyCharacteristicResponse{
			ID:             characteristic.ID,
			Name:           characteristic.Name,
			Many:           characteristic.Many,
			PossibleValues: strings.Split(characteristic.PossibleValues, ","),
		}

		characteristics = append(characteristics, &characteristicResponse)
	}

	return characteristics, nil

}
