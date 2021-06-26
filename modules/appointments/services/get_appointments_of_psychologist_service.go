package services

import (
	"context"

	"github.com/guicostaarantes/psi-server/modules/appointments/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// GetAppointmentsOfPsychologistService is a service that the psychologist will use to retrieve their appointments
type GetAppointmentsOfPsychologistService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s GetAppointmentsOfPsychologistService) Execute(psychologistID string) ([]*models.Appointment, error) {

	appointments := []*models.Appointment{}

	cursor, findErr := s.DatabaseUtil.FindMany("psi_db", "appointments", map[string]interface{}{"psychologistId": psychologistID})
	if findErr != nil {
		return nil, findErr
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		appointment := models.Appointment{}

		decodeErr := cursor.Decode(&appointment)
		if decodeErr != nil {
			return nil, decodeErr
		}

		appointments = append(appointments, &appointment)

	}

	return appointments, nil

}
