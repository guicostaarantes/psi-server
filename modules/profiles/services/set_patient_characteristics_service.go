package services

import (
	"strings"

	models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/identifier"
)

// SetPatientCharacteristicsService is a service that sets all possible patient characteristics
type SetPatientCharacteristicsService struct {
	DatabaseUtil   database.IDatabaseUtil
	IdentifierUtil identifier.IIdentifierUtil
}

// Execute is the method that runs the business logic of the service
func (s SetPatientCharacteristicsService) Execute(input []*models.SetPatientCharacteristicInput) error {

	newCharacteristics := []interface{}{}

	for _, char := range input {
		characteristic := models.PatientCharacteristic{
			Name:           char.Name,
			Many:           char.Many,
			PossibleValues: strings.Join(char.PossibleValues, ","),
		}

		newCharacteristics = append(newCharacteristics, characteristic)
	}

	deleteErr := s.DatabaseUtil.DeleteMany("psi_db", "patient_characteristics", map[string]interface{}{})
	if deleteErr != nil {
		return deleteErr
	}

	writeErr := s.DatabaseUtil.InsertMany("psi_db", "patient_characteristics", newCharacteristics)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
