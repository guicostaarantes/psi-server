package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/guicostaarantes/psi-server/modules/characteristics/models"
	profiles_models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// SetCharacteristicChoicesService is a service that assigns a characteristic to a patient profile
type SetCharacteristicChoicesService struct {
	OrmUtil orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s SetCharacteristicChoicesService) Execute(id string, input []*models.SetCharacteristicChoiceInput) error {

	var target models.CharacteristicTarget

	psy := profiles_models.Psychologist{}
	pat := profiles_models.Patient{}
	result := s.OrmUtil.Db().Where("id = ?", id).Limit(1).Find(&pat)
	if result.Error != nil {
		return result.Error
	}
	if pat.ID != "" {
		target = models.PatientTarget
	} else {
		result := s.OrmUtil.Db().Where("id = ?", id).Limit(1).Find(&psy)
		if result.Error != nil {
			return result.Error
		}
		if psy.ID != "" {
			target = models.PsychologistTarget
		} else {
			return errors.New("resource not found")
		}
	}

	characteristics := []*models.Characteristic{}

	result = s.OrmUtil.Db().Where("target = ?", target).Find(&characteristics)
	if result.Error != nil {
		return result.Error
	}

	characteristicsTypes := map[string]models.CharacteristicType{}
	possibleValues := map[string]map[string]bool{}

	for _, char := range characteristics {
		characteristicsTypes[char.Name] = char.Type
		for _, pv := range strings.Split(char.PossibleValues, ",") {
			if _, exists := possibleValues[char.Name]; !exists {
				possibleValues[char.Name] = map[string]bool{}
			}
			possibleValues[char.Name][pv] = true
		}
	}

	choicesToCreate := []*models.CharacteristicChoice{}

	for _, newChoices := range input {

		switch characteristicsTypes[newChoices.CharacteristicName] {

		case models.Boolean:
			if len(newChoices.SelectedValues) != 1 || (newChoices.SelectedValues[0] != "true" && newChoices.SelectedValues[0] != "false") {
				return fmt.Errorf("characteristic '%s' must be either true or false", newChoices.CharacteristicName)
			}
			choicesToCreate = append(choicesToCreate, &models.CharacteristicChoice{
				ProfileID:          id,
				Target:             target,
				CharacteristicName: newChoices.CharacteristicName,
				SelectedValue:      newChoices.SelectedValues[0],
			})

		case models.Single:
			if len(newChoices.SelectedValues) != 1 {
				return fmt.Errorf("characteristic '%s' needs exactly one value", newChoices.CharacteristicName)
			}
			if _, exists := possibleValues[newChoices.CharacteristicName][newChoices.SelectedValues[0]]; !exists {
				return fmt.Errorf("option '%s' is not possible in characteristic %s", newChoices.SelectedValues[0], newChoices.CharacteristicName)
			}
			choicesToCreate = append(choicesToCreate, &models.CharacteristicChoice{
				ProfileID:          id,
				Target:             target,
				CharacteristicName: newChoices.CharacteristicName,
				SelectedValue:      newChoices.SelectedValues[0],
			})

		case models.Multiple:
			for _, sv := range newChoices.SelectedValues {
				if _, exists := possibleValues[newChoices.CharacteristicName][sv]; !exists {
					return fmt.Errorf("option %s is not possible in characteristic %s", sv, newChoices.CharacteristicName)
				}
				choicesToCreate = append(choicesToCreate, &models.CharacteristicChoice{
					ProfileID:          id,
					Target:             target,
					CharacteristicName: newChoices.CharacteristicName,
					SelectedValue:      sv,
				})
			}

		default:
			return fmt.Errorf("characteristic has unknown type %s", newChoices.CharacteristicName)

		}

	}

	result = s.OrmUtil.Db().Delete(&models.CharacteristicChoice{}, "profile_id = ?", id)
	if result.Error != nil {
		return result.Error
	}

	if len(choicesToCreate) > 0 {
		result = s.OrmUtil.Db().Create(&choicesToCreate)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil

}
