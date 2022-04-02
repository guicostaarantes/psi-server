package appointments_services

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"os"
	"time"

	appointments_models "github.com/guicostaarantes/psi-server/modules/appointments/models"
	appointments_templates "github.com/guicostaarantes/psi-server/modules/appointments/templates"
	mails_models "github.com/guicostaarantes/psi-server/modules/mails/models"
	profiles_models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	users_models "github.com/guicostaarantes/psi-server/modules/users/models"
	"github.com/guicostaarantes/psi-server/utils/identifier"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// EditAppointmentByPsychologistService is a service that the psychologist will use to edit an appointment
type EditAppointmentByPsychologistService struct {
	IdentifierUtil identifier.IIdentifierUtil
	OrmUtil        orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s EditAppointmentByPsychologistService) Execute(id string, psychologistID string, input appointments_models.EditAppointmentByPsychologistInput) error {

	appointment := appointments_models.Appointment{}
	psychologist := profiles_models.Psychologist{}
	patient := profiles_models.Patient{}
	patientUser := users_models.User{}

	result := s.OrmUtil.Db().Where("id = ? AND psychologist_id = ?", id, psychologistID).Limit(1).Find(&appointment)
	if result.Error != nil {
		return result.Error
	}

	if appointment.ID == "" {
		return errors.New("resource not found")
	}

	result = s.OrmUtil.Db().Where("id = ?", psychologistID).Limit(1).Find(&psychologist)
	if result.Error != nil {
		return result.Error
	}

	result = s.OrmUtil.Db().Where("id = ?", appointment.PatientID).Limit(1).Find(&patient)
	if result.Error != nil {
		return result.Error
	}

	result = s.OrmUtil.Db().Where("id = ?", patient.UserID).Limit(1).Find(&patientUser)
	if result.Error != nil {
		return result.Error
	}

	if appointment.Status == appointments_models.CanceledByPatient {
		return fmt.Errorf("appointment status cannot change from %s to EDITED_BY_PSYCHOLOGIST", string(appointment.Status))
	}

	if time.Now().After(input.Start) {
		return errors.New("appointment cannot be scheduled to the past")
	}

	if !input.End.After(input.Start) {
		return errors.New("appointment cannot have negative duration")
	}

	appointment.Status = appointments_models.EditedByPsychologist
	appointment.Start = input.Start
	appointment.End = input.End
	appointment.PriceRangeName = input.PriceRangeName
	appointment.Reason = input.Reason

	_, mailID, mailIDErr := s.IdentifierUtil.GenerateIdentifier()
	if mailIDErr != nil {
		return mailIDErr
	}

	templ, templErr := template.New("AppointmentModifiedByPsychologistEmail").Parse(appointments_templates.AppointmentModifiedByPsychologistEmailTemplate)
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
		Subject:     "Consulta modificada no PSI",
		Html:        buff.String(),
		Processed:   false,
	}

	result = s.OrmUtil.Db().Create(&mail)
	if result.Error != nil {
		return result.Error
	}

	result = s.OrmUtil.Db().Save(&appointment)
	if result.Error != nil {
		return result.Error
	}

	return nil

}
