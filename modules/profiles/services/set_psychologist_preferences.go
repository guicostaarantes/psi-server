package services

import (
	"context"
	"strings"

	"github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// SetPsychologistPreferencesService is a service that allows the psychologist to submit their preferences
type SetPsychologistPreferencesService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s SetPsychologistPreferencesService) Execute(id string, input []*models.SetPsychologistPreferenceInput) error {

	preferences := []interface{}{}

	cursor, findErr := s.DatabaseUtil.FindMany("psi_db", "patient_characteristics", map[string]interface{}{})
	if findErr != nil {
		return findErr
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {

		characteristic := models.PatientCharacteristic{}

		decodeErr := cursor.Decode(&characteristic)
		if decodeErr != nil {
			return decodeErr
		}

		for _, i := range input {
			if characteristic.Name == i.CharacteristicName {
				possibleValues := strings.Split(characteristic.PossibleValues, ",")
				for _, p := range possibleValues {
					if i.Value == p && i.Weight != 0 {
						preferences = append(preferences, models.PsychologistPreference{
							PsychologistID:     id,
							CharacteristicName: i.CharacteristicName,
							Value:              i.Value,
							Weight:             i.Weight,
						})
					}
				}
			}
		}

	}

	deleteErr := s.DatabaseUtil.DeleteMany("psi_db", "psychologist_preferences", map[string]interface{}{"psychologistId": id})
	if deleteErr != nil {
		return deleteErr
	}

	writeErr := s.DatabaseUtil.InsertMany("psi_db", "psychologist_preferences", preferences)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
