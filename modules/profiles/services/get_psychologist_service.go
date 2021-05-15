package services

import (
	models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// GetPsychologistService is a service that gets the psychologist profile based on id
type GetPsychologistService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s GetPsychologistService) Execute(id string) (*models.Psychologist, error) {

	psy := &models.Psychologist{}

	findErr := s.DatabaseUtil.FindOne("psi_db", "psychologists", map[string]interface{}{"id": id}, psy)
	if findErr != nil {
		return nil, findErr
	}

	if psy.ID == "" {
		return nil, nil
	}
	
	return psy, nil

}
