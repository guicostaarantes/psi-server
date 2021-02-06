package services

import (
	"errors"

	models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// UpdatePsychologistService is a service that edits a psychologist profile
type UpdatePsychologistService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s UpdatePsychologistService) Execute(id string, psyInput *models.UpdatePsychologistInput) error {

	psy := models.Psychologist{}

	findErr := s.DatabaseUtil.FindOne("psi_db", "psychologists", map[string]interface{}{"id": id}, &psy)
	if findErr != nil {
		return findErr
	}

	if psy.ID == "" {
		return errors.New("resource not found")
	}

	psy.BirthDate = psyInput.BirthDate
	psy.City = psyInput.City

	writeErr := s.DatabaseUtil.UpdateOne("psi_db", "psychologists", map[string]interface{}{"id": id}, psy)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
