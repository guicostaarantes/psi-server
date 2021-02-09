package services

import (
	"context"
	"strings"

	models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// GetPatientCharacteristicsByPatientIDService is a service that gets the characteristics of a patient by patientId
type GetPatientCharacteristicsByPatientIDService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s GetPatientCharacteristicsByPatientIDService) Execute(id string) ([]*models.PatientCharacteristicChoiceResponse, error) {

	characteristics := []*models.PatientCharacteristicChoiceResponse{}

	charCursor, findErr := s.DatabaseUtil.FindMany("psi_db", "patient_characteristics", map[string]interface{}{})
	if findErr != nil {
		return nil, findErr
	}

	defer charCursor.Close(context.Background())

	for charCursor.Next(context.Background()) {
		characteristic := models.PatientCharacteristic{}

		decodeErr := charCursor.Decode(&characteristic)
		if decodeErr != nil {
			return nil, decodeErr
		}

		characteristicResponse := models.PatientCharacteristicChoiceResponse{
			Name:           characteristic.Name,
			Many:           characteristic.Many,
			Values:         []string{},
			PossibleValues: strings.Split(characteristic.PossibleValues, ","),
		}

		characteristics = append(characteristics, &characteristicResponse)
	}

	choiceCursor, findErr := s.DatabaseUtil.FindMany("psi_db", "patient_characteristic_choices", map[string]interface{}{"patientId": id})
	if findErr != nil {
		return nil, findErr
	}

	defer choiceCursor.Close(context.Background())

	for choiceCursor.Next(context.Background()) {
		choice := models.PatientCharacteristicChoice{}

		decodeErr := choiceCursor.Decode(&choice)
		if decodeErr != nil {
			return nil, decodeErr
		}

		for _, char := range characteristics {
			if char.Name == choice.CharacteristicName {
				char.Values = append(char.Values, choice.Value)
			}
		}
	}

	return characteristics, nil

}
