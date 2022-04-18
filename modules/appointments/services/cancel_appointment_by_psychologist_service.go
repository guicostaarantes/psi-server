package appointments_services

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"os"

	appointments_models "github.com/guicostaarantes/psi-server/modules/appointments/models"
	appointments_templates "github.com/guicostaarantes/psi-server/modules/appointments/templates"
	mails_models "github.com/guicostaarantes/psi-server/modules/mails/models"
	profiles_models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	users_models "github.com/guicostaarantes/psi-server/modules/users/models"
	"github.com/guicostaarantes/psi-server/utils/identifier"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// CancelAppointmentByPsychologistService is a service that the psychologist will use to cancel an appointment
type CancelAppointmentByPsychologistService struct {
	IdentifierUtil identifier.IIdentifierUtil
	OrmUtil        orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s CancelAppointmentByPsychologistService) Execute(id string, psychologistID string, reason string) error {

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

	if appointment.Status == appointments_models.CanceledByPatient || appointment.Status == appointments_models.CanceledByPsychologist {
		return fmt.Errorf("appointment status cannot change from %s to CANCELED_BY_PSYCHOLOGIST", string(appointment.Status))
	}

	appointment.Status = appointments_models.CanceledByPsychologist
	appointment.Reason = reason

	_, mailID, mailIDErr := s.IdentifierUtil.GenerateIdentifier()
	if mailIDErr != nil {
		return mailIDErr
	}

	templ, templErr := template.New("AppointmentCanceledByPsychologistEmail").Parse(appointments_templates.AppointmentCanceledByPsychologistEmailTemplate)
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
		Subject:     "Consulta cancelada no PSI",
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
