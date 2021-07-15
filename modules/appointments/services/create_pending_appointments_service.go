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
	DatabaseUtil   database.IDatabaseUtil
	IdentifierUtil identifier.IIdentifierUtil
}

// Execute is the method that runs the business logic of the service
func (s CreatePendingAppointmentsService) Execute() error {

	activeTreatments := map[string]*treatments_models.Treatment{}

	cursor, findErr := s.DatabaseUtil.FindMany("psi_db", "treatments", map[string]interface{}{"status": string(treatments_models.Active)})
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

	cursor, findErr = s.DatabaseUtil.FindMany("psi_db", "appointments", map[string]interface{}{})
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
		// logic to determine next appointment's start based on treatment.WeeklyStart
		currentTime := time.Now().Unix()
		secondsInOneWeek := int64(7 * 24 * 60 * 60)
		phaseAdjustment := int64(3 * 24 * 60 * 60) // unix time 0 is a Thursday, but out weekly calendar starts on a Sunday
		startOfCurrentWeek := (currentTime-phaseAdjustment)/secondsInOneWeek*secondsInOneWeek + phaseAdjustment
		start := startOfCurrentWeek + treatment.WeeklyStart

		// if the start time of the current week has already passed, send it to the week after
		if currentTime-startOfCurrentWeek >= treatment.WeeklyStart {
			start += secondsInOneWeek
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
			Start:          start,
			End:            start + treatment.Duration,
			Price:          treatment.Price,
			Status:         models.Created,
		}

		appointmentsToCreate = append(appointmentsToCreate, newAppointment)
	}

	writeErr := s.DatabaseUtil.InsertMany("psi_db", "appointments", appointmentsToCreate)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
