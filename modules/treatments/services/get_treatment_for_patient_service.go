package services

import (
	"github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// GetTreatmentForPatientService is a service that gets a treatment based on its id
type GetTreatmentForPatientService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s GetTreatmentForPatientService) Execute(id string) (*models.GetPatientTreatmentsResponse, error) {

	treatment := &models.GetPatientTreatmentsResponse{}

	findErr := s.DatabaseUtil.FindOne("psi_db", "treatments", map[string]interface{}{"id": id}, &treatment)
	if findErr != nil {
		return nil, findErr
	}

	return treatment, nil

}
