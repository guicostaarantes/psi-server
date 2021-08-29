package services

import (
	"context"
	"time"

	"github.com/guicostaarantes/psi-server/modules/appointments/models"
	treatments_models "github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/identifier"
)

// CreatePendingAppointmentsService is a service that creates appointments for all active treatments that have no appointments scheduled to the future
type CreatePendingAppointmentsService struct {
	DatabaseUtil            database.IDatabaseUtil
	IdentifierUtil          identifier.IIdentifierUtil
	ScheduleIntervalSeconds int64
}

// Execute is the method that runs the business logic of the service
func (s CreatePendingAppointmentsService) Execute() error {

	activeTreatments := map[string]*treatments_models.Treatment{}

	cursor, findErr := s.DatabaseUtil.FindMany("treatments", map[string]interface{}{"status": string(treatments_models.Active)})
	if findErr != nil {
		return findErr
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		treatment := treatments_models.Treatment{}

		decodeErr := cursor.Decode(&treatment)
		if decodeErr != nil {
			return decodeErr
		}

		activeTreatments[treatment.ID] = &treatment
	}

	cursor, findErr = s.DatabaseUtil.FindMany("appointments", map[string]interface{}{})
	if findErr != nil {
		return findErr
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		appointment := models.Appointment{}

		decodeErr := cursor.Decode(&appointment)
		if decodeErr != nil {
			return decodeErr
		}

		_, exists := activeTreatments[appointment.TreatmentID]

		if exists && appointment.End > time.Now().Unix() {
			delete(activeTreatments, appointment.TreatmentID)
		}
	}

	appointmentsToCreate := []interface{}{}

	for id, treatment := range activeTreatments {
		currentTime := time.Now().Unix()
		intervalDuration := s.ScheduleIntervalSeconds * treatment.Frequency
		currentInterval := currentTime / intervalDuration
		nextAppointmentStart := intervalDuration*currentInterval + treatment.Phase
		// if the start time of the current interval has already passed, send it to the next interval
		if nextAppointmentStart <= currentTime {
			nextAppointmentStart += intervalDuration
		}

		_, appoID, appoIDErr := s.IdentifierUtil.GenerateIdentifier()
		if appoIDErr != nil {
			return appoIDErr
		}

		newAppointment := models.Appointment{
			ID:             appoID,
			TreatmentID:    id,
			PatientID:      treatment.PatientID,
			PsychologistID: treatment.PsychologistID,
			Start:          nextAppointmentStart,
			End:            nextAppointmentStart + treatment.Duration,
			PriceRange:     treatment.PriceRange,
			Status:         models.Created,
		}

		appointmentsToCreate = append(appointmentsToCreate, newAppointment)
	}

	writeErr := s.DatabaseUtil.InsertMany("appointments", appointmentsToCreate)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
