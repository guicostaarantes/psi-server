package profiles_services

import (
	"context"
	"errors"
	"strings"

	models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/merge"
)

// SetPsyCharacteristicChoiceService is a service that assigns a characteristic to a psychologist profile
type SetPsyCharacteristicChoiceService struct {
	DatabaseUtil database.IDatabaseUtil
	MergeUtil    merge.IMergeUtil
}

// Execute is the method that runs the business logic of the service
func (s SetPsyCharacteristicChoiceService) Execute(psyChoiceInput *models.SetPsyCharacteristicChoiceInput) error {

	characteristic := models.PsyCharacteristic{}

	findErr := s.DatabaseUtil.FindOne("psi_db", "psychologist_characteristics", map[string]interface{}{"name": psyChoiceInput.CharacteristicName}, &characteristic)
	if findErr != nil {
		return findErr
	}

	if !characteristic.Many && len(psyChoiceInput.Values) != 1 {
		return errors.New("characteristic '" + characteristic.Name + "' needs exactly one value")
	}

	possibleValues := strings.Split(characteristic.PossibleValues, ",")

	for _, value := range psyChoiceInput.Values {
		var valueExists = false
		for _, posValue := range possibleValues {
			if posValue == value {
				valueExists = true
			}
		}
		if !valueExists {
			return errors.New("option '" + value + "' is not possible")
		}
	}

	otherChoicesCursor, findErr := s.DatabaseUtil.FindMany("psi_db", "psychologist_characteristic_choices", map[string]interface{}{
		"psychologistId":     psyChoiceInput.PsychologistID,
		"characteristicName": psyChoiceInput.CharacteristicName,
	})
	if findErr != nil {
		return findErr
	}

	defer otherChoicesCursor.Close(context.Background())

	for otherChoicesCursor.Next(context.Background()) {
		choice := models.PsyCharacteristicChoice{}

		decodeErr := otherChoicesCursor.Decode(&choice)
		if decodeErr != nil {
			return decodeErr
		}

		deleteErr := s.DatabaseUtil.DeleteOne("psi_db", "psychologist_characteristic_choices", map[string]interface{}{
			"psychologistId":     choice.PsychologistID,
			"characteristicName": choice.CharacteristicName,
		})
		if deleteErr != nil {
			return deleteErr
		}
	}

	for _, value := range psyChoiceInput.Values {
		choice := models.PsyCharacteristicChoice{
			PsychologistID:     psyChoiceInput.PsychologistID,
			CharacteristicName: psyChoiceInput.CharacteristicName,
			Value:              value,
		}

		insertErr := s.DatabaseUtil.InsertOne("psi_db", "psychologist_characteristic_choices", choice)
		if insertErr != nil {
			return insertErr
		}
	}

	return nil

}
