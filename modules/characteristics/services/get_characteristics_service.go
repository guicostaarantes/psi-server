package characteristcs_services

import (
	"strings"

	characteristics_models "github.com/guicostaarantes/psi-server/modules/characteristics/models"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// GetCharacteristicsService is a service that gets all possible characteristic based on the target
type GetCharacteristicsService struct {
	OrmUtil orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s GetCharacteristicsService) Execute(target characteristics_models.CharacteristicTarget) ([]*characteristics_models.CharacteristicResponse, error) {

	response := []*characteristics_models.CharacteristicResponse{}
	characteristics := []*characteristics_models.Characteristic{}

	result := s.OrmUtil.Db().Where("target = ?", target).Find(&characteristics)
	if result.Error != nil {
		return nil, result.Error
	}

	for _, char := range characteristics {
		response = append(response, &characteristics_models.CharacteristicResponse{
			Name:           char.Name,
			Type:           char.Type,
			PossibleValues: strings.Split(char.PossibleValues, ","),
		})
	}

	return response, nil

}
