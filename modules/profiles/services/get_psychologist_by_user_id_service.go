package services

import (
	models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// GetPsychologistByUserIDService is a service that gets the psychologist profile based on UserID
type GetPsychologistByUserIDService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s GetPsychologistByUserIDService) Execute(id string) (*models.Psychologist, error) {

	psy := &models.Psychologist{}

	findErr := s.DatabaseUtil.FindOne("psychologists", map[string]interface{}{"userId": id}, psy)
	if findErr != nil {
		return nil, findErr
	}

	if psy.ID == "" {
		return nil, nil
	}

	return psy, nil

}
