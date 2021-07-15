package services

import (
	"context"

	"github.com/guicostaarantes/psi-server/modules/appointments/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// GetAppointmentsOfPatientService is a service that the patient will use to retrieve their appointments
type GetAppointmentsOfPatientService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s GetAppointmentsOfPatientService) Execute(patientID string) ([]*models.Appointment, error) {

	appointments := []*models.Appointment{}

	cursor, findErr := s.DatabaseUtil.FindMany("psi_db", "appointments", map[string]interface{}{"patientId": patientID})
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
