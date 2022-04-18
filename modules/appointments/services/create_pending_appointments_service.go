package appointments_services

import (
	"bytes"
	"html/template"
	"os"
	"time"

	appointments_models "github.com/guicostaarantes/psi-server/modules/appointments/models"
	appointments_templates "github.com/guicostaarantes/psi-server/modules/appointments/templates"
	mails_models "github.com/guicostaarantes/psi-server/modules/mails/models"
	profiles_models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	treatments_models "github.com/guicostaarantes/psi-server/modules/treatments/models"
	users_models "github.com/guicostaarantes/psi-server/modules/users/models"
	"github.com/guicostaarantes/psi-server/utils/identifier"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// CreatePendingAppointmentsService is a service that creates appointments for all active treatments that have no appointments scheduled to the future
type CreatePendingAppointmentsService struct {
	IdentifierUtil           identifier.IIdentifierUtil
	OrmUtil                  orm.IOrmUtil
	ScheduleIntervalDuration time.Duration
}

// Execute is the method that runs the business logic of the service
func (s CreatePendingAppointmentsService) Execute() error {

	activeTreatmentsWithoutFutureAppointments := []*treatments_models.Treatment{}

	result := s.OrmUtil.Db().Raw(
		"SELECT * FROM treatments WHERE id IN (SELECT DISTINCT treatments.id FROM treatments LEFT JOIN appointments ON appointments.treatment_id = treatments.id WHERE treatments.status = ? EXCEPT SELECT treatment_id FROM appointments WHERE start > ?)",
		treatments_models.Active,
		time.Now(),
	).Find(&activeTreatmentsWithoutFutureAppointments)
	if result.Error != nil {
		return result.Error
	}

	for _, treatment := range activeTreatmentsWithoutFutureAppointments {
		currentTime := time.Now()
		intervalDuration := int64(s.ScheduleIntervalDuration/time.Second) * treatment.Frequency
		currentInterval := currentTime.Unix() / intervalDuration
		nextAppointmentStart := time.Unix(intervalDuration*currentInterval+treatment.Phase, 10)
		// if the start time of the current interval has already passed, send it to the next interval
		if currentTime.After(nextAppointmentStart) {
			nextAppointmentStart = nextAppointmentStart.Add(time.Duration(intervalDuration) * time.Second)
		}

		_, appoID, appoIDErr := s.IdentifierUtil.GenerateIdentifier()
		if appoIDErr != nil {
			return appoIDErr
		}

		newAppointment := appointments_models.Appointment{
			ID:             appoID,
			TreatmentID:    treatment.ID,
			PatientID:      treatment.PatientID,
			PsychologistID: treatment.PsychologistID,
			Start:          nextAppointmentStart,
			End:            nextAppointmentStart.Add(time.Duration(treatment.Duration) * time.Second),
			PriceRangeName: treatment.PriceRangeName,
			Status:         appointments_models.Created,
		}

		psychologist := profiles_models.Psychologist{}
		patient := profiles_models.Patient{}
		patientUser := users_models.User{}

		result = s.OrmUtil.Db().Where("id = ?", treatment.PsychologistID).Limit(1).Find(&psychologist)
		if result.Error != nil {
			return result.Error
		}

		result = s.OrmUtil.Db().Where("id = ?", treatment.PatientID).Limit(1).Find(&patient)
		if result.Error != nil {
			return result.Error
		}

		result = s.OrmUtil.Db().Where("id = ?", patient.UserID).Limit(1).Find(&patientUser)
		if result.Error != nil {
			return result.Error
		}

		_, mailID, mailIDErr := s.IdentifierUtil.GenerateIdentifier()
		if mailIDErr != nil {
			return mailIDErr
		}

		templ, templErr := template.New("AppointmentCreatedEmail").Parse(appointments_templates.AppointmentCreatedEmailTemplate)
		if templErr != nil {
			return templErr
		}

		buff := new(bytes.Buffer)

		templ.Execute(buff, map[string]string{
			"SiteURL":     os.Getenv("PSI_SITE_URL"),
			"LikeName":    patient.LikeName,
			"PsyFullName": psychologist.FullName,
		})

		mail := &mails_models.TransientMailMessage{
			ID:          mailID,
			FromAddress: "relacionamento@psi.com.br",
			FromName:    "Relacionamento PSI",
			To:          patientUser.Email,
			Cc:          "",
			Cco:         "",
			Subject:     "Consulta criada no PSI",
			Html:        buff.String(),
			Processed:   false,
		}

		result = s.OrmUtil.Db().Create(&mail)
		if result.Error != nil {
			return result.Error
		}

		result = s.OrmUtil.Db().Create(&newAppointment)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil

}
