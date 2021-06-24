package services

import (
	"github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// GetTreatmentForPsychologistService is a service that gets a treatment based on its id
type GetTreatmentForPsychologistService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s GetTreatmentForPsychologistService) Execute(id string) (*models.GetPsychologistTreatmentsResponse, error) {

	treatment := &models.GetPsychologistTreatmentsResponse{}

	findErr := s.DatabaseUtil.FindOne("psi_db", "treatments", map[string]interface{}{"id": id}, &treatment)
	if findErr != nil {
		return nil, findErr
	}

	return treatment, nil

}
