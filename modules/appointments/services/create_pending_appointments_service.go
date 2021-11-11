package services

import (
	"time"

	"github.com/guicostaarantes/psi-server/modules/appointments/models"
	treatments_models "github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/identifier"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// CreatePendingAppointmentsService is a service that creates appointments for all active treatments that have no appointments scheduled to the future
type CreatePendingAppointmentsService struct {
	IdentifierUtil          identifier.IIdentifierUtil
	OrmUtil                 orm.IOrmUtil
	ScheduleIntervalSeconds int64
}

// Execute is the method that runs the business logic of the service
func (s CreatePendingAppointmentsService) Execute() error {

	activeTreatmentsWithoutFutureAppointments := []*treatments_models.Treatment{}

	result := s.OrmUtil.Db().Raw(
		"SELECT * FROM treatments WHERE id IN (SELECT DISTINCT treatments.id FROM treatments LEFT JOIN appointments ON appointments.treatment_id = treatments.id WHERE treatments.status = ? EXCEPT SELECT treatment_id FROM appointments WHERE start > ?)",
		treatments_models.Active,
		time.Now().Unix(),
	).Find(&activeTreatmentsWithoutFutureAppointments)
	if result.Error != nil {
		return result.Error
	}

	appointmentsToCreate := []*models.Appointment{}

	for _, treatment := range activeTreatmentsWithoutFutureAppointments {
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
			TreatmentID:    treatment.ID,
			PatientID:      treatment.PatientID,
			PsychologistID: treatment.PsychologistID,
			Start:          nextAppointmentStart,
			End:            nextAppointmentStart + treatment.Duration,
			PriceRangeName: treatment.PriceRangeName,
			Status:         models.Created,
		}

		appointmentsToCreate = append(appointmentsToCreate, &newAppointment)
	}

	if len(appointmentsToCreate) > 0 {
		result = s.OrmUtil.Db().Create(&appointmentsToCreate)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil

}
