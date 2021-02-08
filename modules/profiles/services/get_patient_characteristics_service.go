package services

import (
	"context"
	"strings"

	models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// GetPatientCharacteristicsService is a service that gets the user by userId
type GetPatientCharacteristicsService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s GetPatientCharacteristicsService) Execute() ([]*models.PatientCharacteristicResponse, error) {

	characteristics := []*models.PatientCharacteristicResponse{}

	cursor, findErr := s.DatabaseUtil.FindMany("psi_db", "patient_characteristics", map[string]interface{}{})
	if findErr != nil {
		return nil, findErr
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		characteristic := models.PatientCharacteristic{}

		decodeErr := cursor.Decode(&characteristic)
		if decodeErr != nil {
			return nil, decodeErr
		}

		characteristicResponse := models.PatientCharacteristicResponse{
			Name:           characteristic.Name,
			Many:           characteristic.Many,
			PossibleValues: strings.Split(characteristic.PossibleValues, ","),
		}

		characteristics = append(characteristics, &characteristicResponse)
	}

	return characteristics, nil

}
