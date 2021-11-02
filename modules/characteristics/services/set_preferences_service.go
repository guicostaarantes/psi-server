package services

import (
	"context"
	"errors"
	"strings"

	"github.com/guicostaarantes/psi-server/modules/characteristics/models"
	profiles_models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// SetPreferencesService is a service that allows a profile to submit their preferences
type SetPreferencesService struct {
	DatabaseUtil database.IDatabaseUtil
	OrmUtil      orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s SetPreferencesService) Execute(id string, input []*models.SetPreferenceInput) error {

	var target models.CharacteristicTarget
	var profileType models.CharacteristicTarget

	psy := profiles_models.Psychologist{}
	pat := profiles_models.Patient{}
	result := s.OrmUtil.Db().Where("id = ?", id).Limit(1).Find(&pat)
	if result.Error != nil {
		return result.Error
	}
	if pat.ID != "" {
		target = models.PsychologistTarget
		profileType = models.PatientTarget
	} else {
		result := s.OrmUtil.Db().Where("id = ?", id).Limit(1).Find(&psy)
		if result.Error != nil {
			return result.Error
		}
		if psy.ID != "" {
			target = models.PatientTarget
			profileType = models.PsychologistTarget
		} else {
			return errors.New("resource not found")
		}
	}

	preferences := []interface{}{}

	cursor, findErr := s.DatabaseUtil.FindMany("characteristics", map[string]interface{}{"target": string(target)})
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
					if i.SelectedValue == p {
						preferences = append(preferences, models.Preference{
							ProfileID:          id,
							Target:             profileType,
							CharacteristicName: i.CharacteristicName,
							SelectedValue:      i.SelectedValue,
							Weight:             i.Weight,
						})
					}
				}
			}
		}

	}

	deleteErr := s.DatabaseUtil.DeleteMany("preferences", map[string]interface{}{"profileId": id})
	if deleteErr != nil {
		return deleteErr
	}

	writeErr := s.DatabaseUtil.InsertMany("preferences", preferences)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
