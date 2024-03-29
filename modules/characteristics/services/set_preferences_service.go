package characteristcs_services

import (
	"errors"
	"fmt"
	"strings"

	characteristics_models "github.com/guicostaarantes/psi-server/modules/characteristics/models"
	profiles_models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// SetPreferencesService is a service that allows a profile to submit their preferences
type SetPreferencesService struct {
	OrmUtil orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s SetPreferencesService) Execute(id string, input []*characteristics_models.SetPreferenceInput) error {

	var target characteristics_models.CharacteristicTarget
	var profileType characteristics_models.CharacteristicTarget

	psy := profiles_models.Psychologist{}
	pat := profiles_models.Patient{}
	result := s.OrmUtil.Db().Where("id = ?", id).Limit(1).Find(&pat)
	if result.Error != nil {
		return result.Error
	}
	if pat.ID != "" {
		target = characteristics_models.PsychologistTarget
		profileType = characteristics_models.PatientTarget
	} else {
		result := s.OrmUtil.Db().Where("id = ?", id).Limit(1).Find(&psy)
		if result.Error != nil {
			return result.Error
		}
		if psy.ID != "" {
			target = characteristics_models.PatientTarget
			profileType = characteristics_models.PsychologistTarget
		} else {
			return errors.New("resource not found")
		}
	}

	characteristics := []*characteristics_models.Characteristic{}

	result = s.OrmUtil.Db().Where("target = ?", target).Find(&characteristics)
	if result.Error != nil {
		return result.Error
	}

	possibleValues := map[string]map[string]bool{}

	for _, char := range characteristics {
		for _, pv := range strings.Split(char.PossibleValues, ",") {
			if _, exists := possibleValues[char.Name]; !exists {
				possibleValues[char.Name] = map[string]bool{}
			}
			possibleValues[char.Name][pv] = true
		}
	}

	preferencesToCreate := []*characteristics_models.Preference{}

	for _, i := range input {
		if _, exists := possibleValues[i.CharacteristicName][i.SelectedValue]; !exists {
			return fmt.Errorf("option '%s' is not possible in characteristic %s", i.SelectedValue, i.CharacteristicName)
		}
		preferencesToCreate = append(preferencesToCreate, &characteristics_models.Preference{
			ProfileID:          id,
			Target:             profileType,
			CharacteristicName: i.CharacteristicName,
			SelectedValue:      i.SelectedValue,
			Weight:             i.Weight,
		})
	}

	result = s.OrmUtil.Db().Delete(&characteristics_models.Preference{}, "profile_id = ?", id)
	if result.Error != nil {
		return result.Error
	}

	if len(preferencesToCreate) > 0 {
		result = s.OrmUtil.Db().Create(&preferencesToCreate)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil

}
