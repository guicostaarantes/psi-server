package characteristcs_services

import (
	"errors"
	"strings"

	characteristics_models "github.com/guicostaarantes/psi-server/modules/characteristics/models"
	profiles_models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// GetCharacteristicsByIDService is a service that gets the characteristics of a profile based on its id
type GetCharacteristicsByIDService struct {
	OrmUtil orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s GetCharacteristicsByIDService) Execute(id string) ([]*characteristics_models.CharacteristicChoiceResponse, error) {

	var target characteristics_models.CharacteristicTarget

	psy := profiles_models.Psychologist{}
	pat := profiles_models.Patient{}
	result := s.OrmUtil.Db().Where("id = ?", id).Limit(1).Find(&pat)
	if result.Error != nil {
		return nil, result.Error
	}
	if pat.ID != "" {
		target = characteristics_models.PatientTarget
	} else {
		result := s.OrmUtil.Db().Where("id = ?", id).Limit(1).Find(&psy)
		if result.Error != nil {
			return nil, result.Error
		}
		if psy.ID != "" {
			target = characteristics_models.PsychologistTarget
		} else {
			return nil, errors.New("resource not found")
		}
	}

	response := []*characteristics_models.CharacteristicChoiceResponse{}
	characteristics := []*characteristics_models.Characteristic{}
	characteristicsChoices := []*characteristics_models.CharacteristicChoice{}

	result = s.OrmUtil.Db().Where("target = ?", target).Find(&characteristics)
	if result.Error != nil {
		return nil, result.Error
	}

	result = s.OrmUtil.Db().Where("profile_id = ?", id).Find(&characteristicsChoices)
	if result.Error != nil {
		return nil, result.Error
	}

	for _, char := range characteristics {
		response = append(response, &characteristics_models.CharacteristicChoiceResponse{
			Name:           char.Name,
			Type:           char.Type,
			SelectedValues: []string{},
			PossibleValues: strings.Split(char.PossibleValues, ","),
		})
	}

	for _, choice := range characteristicsChoices {
		for _, char := range response {
			if char.Name == choice.CharacteristicName {
				char.SelectedValues = append(char.SelectedValues, choice.SelectedValue)
			}
		}
	}

	return response, nil

}
