package services

import (
	"context"
	"errors"
	"strings"

	"github.com/guicostaarantes/psi-server/modules/characteristics/models"
	profiles_models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// SetCharacteristicChoicesService is a service that assigns a characteristic to a patient profile
type SetCharacteristicChoicesService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s SetCharacteristicChoicesService) Execute(id string, input []*models.SetCharacteristicChoiceInput) error {

	var target models.CharacteristicTarget

	psy := profiles_models.Psychologist{}
	pat := profiles_models.Patient{}
	s.DatabaseUtil.FindOne("psi_db", "patients", map[string]interface{}{"id": id}, &pat)
	if pat.ID != "" {
		target = models.PatientTarget
	} else {
		s.DatabaseUtil.FindOne("psi_db", "psychologists", map[string]interface{}{"id": id}, &psy)
		if psy.ID != "" {
			target = models.PsychologistTarget
		} else {
			return errors.New("resource not found")
		}
	}

	choices := []interface{}{}

	cursor, findErr := s.DatabaseUtil.FindMany("psi_db", "characteristics", map[string]interface{}{"target": string(target)})
	if findErr != nil {
		return findErr
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		characteristic := models.Characteristic{}

		decodeErr := cursor.Decode(&characteristic)
		if decodeErr != nil {
			return decodeErr
		}

		for _, i := range input {

			if characteristic.Name == i.CharacteristicName {

				if characteristic.Type == models.Boolean {

					if len(i.SelectedValues) != 1 || (i.SelectedValues[0] != "true" && i.SelectedValues[0] != "false") {
						return errors.New("characteristic '" + characteristic.Name + "' must be either true or false")
					}

					choices = append(choices, models.CharacteristicChoice{
						ProfileID:          id,
						CharacteristicName: i.CharacteristicName,
						SelectedValue:      i.SelectedValues[0],
					})

					continue

				}

				if characteristic.Type == models.Single {

					if len(i.SelectedValues) != 1 {
						return errors.New("characteristic '" + characteristic.Name + "' needs exactly one value")
					}

					possibleValues := strings.Split(characteristic.PossibleValues, ",")

					var valueExists = false
					for _, posValue := range possibleValues {
						if posValue == i.SelectedValues[0] {
							valueExists = true
						}
					}

					if !valueExists {
						return errors.New("option '" + i.SelectedValues[0] + "' is not possible in characteristic " + i.CharacteristicName)
					}

					choices = append(choices, models.CharacteristicChoice{
						ProfileID:          id,
						CharacteristicName: i.CharacteristicName,
						SelectedValue:      i.SelectedValues[0],
					})

					continue

				}

				if characteristic.Type == models.Multiple {

					possibleValues := strings.Split(characteristic.PossibleValues, ",")

					for _, value := range i.SelectedValues {
						var valueExists = false
						for _, posValue := range possibleValues {
							if posValue == value {
								valueExists = true
							}
						}
						if !valueExists {
							return errors.New("option '" + value + "' is not possible in characteristic " + i.CharacteristicName)
						}

						choices = append(choices, models.CharacteristicChoice{
							ProfileID:          id,
							CharacteristicName: i.CharacteristicName,
							SelectedValue:      value,
						})
					}

					continue

				}

			}

		}

	}

	deleteErr := s.DatabaseUtil.DeleteMany("psi_db", "characteristic_choices", map[string]interface{}{"profileId": id})
	if deleteErr != nil {
		return deleteErr
	}

	writeErr := s.DatabaseUtil.InsertMany("psi_db", "characteristic_choices", choices)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
