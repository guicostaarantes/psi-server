package services

import (
	"context"
	"errors"
	"strings"

	models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// SetPsychologistCharacteristicChoicesService is a service that assigns a characteristic to a psychologist profile
type SetPsychologistCharacteristicChoicesService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s SetPsychologistCharacteristicChoicesService) Execute(id string, input []*models.SetPsychologistCharacteristicChoiceInput) error {

	choices := []interface{}{}

	cursor, findErr := s.DatabaseUtil.FindMany("psi_db", "psychologist_characteristics", map[string]interface{}{})
	if findErr != nil {
		return findErr
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		characteristic := models.PsychologistCharacteristic{}

		decodeErr := cursor.Decode(&characteristic)
		if decodeErr != nil {
			return decodeErr
		}

		for _, i := range input {

			if characteristic.Name == i.CharacteristicName {

				if !characteristic.Many && characteristic.PossibleValues == "true,false" {

					if len(i.Values) != 1 || (i.Values[0] != "true" && i.Values[0] != "false") {
						return errors.New("characteristic '" + characteristic.Name + "' must be either true or false")
					}

					if i.Values[0] == "true" {
						choices = append(choices, models.PsychologistCharacteristicChoice{
							PsychologistID:     id,
							CharacteristicName: i.CharacteristicName,
							Value:              "true",
						})
					}

					continue

				}

				if !characteristic.Many && len(i.Values) != 1 {
					return errors.New("characteristic '" + characteristic.Name + "' needs exactly one value")
				}

				possibleValues := strings.Split(characteristic.PossibleValues, ",")

				for _, value := range i.Values {
					var valueExists = false
					for _, posValue := range possibleValues {
						if posValue == value {
							valueExists = true
						}
					}
					if !valueExists {
						return errors.New("option '" + value + "' is not possible in characteristic " + i.CharacteristicName)
					}

					choices = append(choices, models.PsychologistCharacteristicChoice{
						PsychologistID:     id,
						CharacteristicName: i.CharacteristicName,
						Value:              value,
					})
				}

			}

		}

	}

	deleteErr := s.DatabaseUtil.DeleteMany("psi_db", "psychologist_characteristic_choices", map[string]interface{}{"psychologistId": id})
	if deleteErr != nil {
		return deleteErr
	}

	writeErr := s.DatabaseUtil.InsertMany("psi_db", "psychologist_characteristic_choices", choices)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
