package services

import (
	"strings"

	models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// SetPsychologistCharacteristicsService is a service that sets all possible psychologist characteristics
type SetPsychologistCharacteristicsService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s SetPsychologistCharacteristicsService) Execute(input []*models.SetPsychologistCharacteristicInput) error {

	newCharacteristics := []interface{}{}

	for _, char := range input {
		characteristic := models.PsychologistCharacteristic{
			Name:           char.Name,
			Many:           char.Many,
			PossibleValues: strings.Join(char.PossibleValues, ","),
		}

		newCharacteristics = append(newCharacteristics, characteristic)
	}

	deleteErr := s.DatabaseUtil.DeleteMany("psi_db", "psychologist_characteristics", map[string]interface{}{})
	if deleteErr != nil {
		return deleteErr
	}

	writeErr := s.DatabaseUtil.InsertMany("psi_db", "psychologist_characteristics", newCharacteristics)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
