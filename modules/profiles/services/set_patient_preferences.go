package services

import (
	"context"
	"strings"

	"github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// SetPatientPreferencesService is a service that allows the patient to submit their preferences
type SetPatientPreferencesService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s SetPatientPreferencesService) Execute(id string, input []*models.SetPatientPreferenceInput) error {

	preferences := []interface{}{}

	cursor, findErr := s.DatabaseUtil.FindMany("psi_db", "psychologist_characteristics", map[string]interface{}{})
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
						preferences = append(preferences, models.PatientPreference{
							PatientID:          id,
							CharacteristicName: i.CharacteristicName,
							Value:              i.Value,
							Weight:             i.Weight,
						})
					}
				}
			}
		}

	}

	deleteErr := s.DatabaseUtil.DeleteMany("psi_db", "patient_preferences", map[string]interface{}{"patientId": id})
	if deleteErr != nil {
		return deleteErr
	}

	writeErr := s.DatabaseUtil.InsertMany("psi_db", "patient_preferences", preferences)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
