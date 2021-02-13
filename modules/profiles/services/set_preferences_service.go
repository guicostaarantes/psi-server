package services

import (
	"context"
	"errors"
	"strings"

	"github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// SetPreferencesService is a service that allows a profile to submit their preferences
type SetPreferencesService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s SetPreferencesService) Execute(id string, input []*models.SetPreferenceInput) error {

	var target models.CharacteristicTarget

	psy := models.Psychologist{}
	pat := models.Patient{}
	s.DatabaseUtil.FindOne("psi_db", "patients", map[string]interface{}{"id": id}, &pat)
	if pat.ID != "" {
		target = models.PsychologistTarget
	} else {
		s.DatabaseUtil.FindOne("psi_db", "psychologists", map[string]interface{}{"id": id}, &psy)
		if psy.ID != "" {
			target = models.PatientTarget
		} else {
			return errors.New("resource not found")
		}
	}

	preferences := []interface{}{}

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
				possibleValues := strings.Split(characteristic.PossibleValues, ",")
				for _, p := range possibleValues {
					if i.SelectedValue == p && i.Weight != 0 {
						preferences = append(preferences, models.Preference{
							ProfileID:          id,
							CharacteristicName: i.CharacteristicName,
							SelectedValue:      i.SelectedValue,
							Weight:             i.Weight,
						})
					}
				}
			}
		}

	}

	deleteErr := s.DatabaseUtil.DeleteMany("psi_db", "preferences", map[string]interface{}{"profileId": id})
	if deleteErr != nil {
		return deleteErr
	}

	writeErr := s.DatabaseUtil.InsertMany("psi_db", "preferences", preferences)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
