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

// CancelAppointmentByPatientService is a service that the patient will use to cancel an appointment
type CancelAppointmentByPatientService struct {
	IdentifierUtil identifier.IIdentifierUtil
	OrmUtil        orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s CancelAppointmentByPatientService) Execute(id string, patientID string, reason string) error {

	appointment := appointments_models.Appointment{}
	patient := profiles_models.Patient{}
	psychologist := profiles_models.Psychologist{}
	psyUser := users_models.User{}

	result := s.OrmUtil.Db().Where("id = ? AND patient_id = ?", id, patientID).Limit(1).Find(&appointment)
	if result.Error != nil {
		return result.Error
	}

	if appointment.ID == "" {
		return errors.New("resource not found")
	}

	result = s.OrmUtil.Db().Where("id = ?", patientID).Limit(1).Find(&patient)
	if result.Error != nil {
		return result.Error
	}

	result = s.OrmUtil.Db().Where("id = ?", appointment.PsychologistID).Limit(1).Find(&psychologist)
	if result.Error != nil {
		return result.Error
	}

	result = s.OrmUtil.Db().Where("id = ?", psychologist.UserID).Limit(1).Find(&psyUser)
	if result.Error != nil {
		return result.Error
	}

	if appointment.Status == appointments_models.CanceledByPatient || appointment.Status == appointments_models.CanceledByPsychologist {
		return fmt.Errorf("appointment status cannot change from %s to CANCELED_BY_PATIENT", string(appointment.Status))
	}

	appointment.Status = appointments_models.CanceledByPatient
	appointment.Reason = reason

	_, mailID, mailIDErr := s.IdentifierUtil.GenerateIdentifier()
	if mailIDErr != nil {
		return mailIDErr
	}

	templ, templErr := template.New("AppointmentCanceledByPatientEmail").Parse(appointments_templates.AppointmentCanceledByPatientEmailTemplate)
	if templErr != nil {
		return templErr
	}

	buff := new(bytes.Buffer)

	templ.Execute(buff, map[string]string{
		"SiteURL":         os.Getenv("PSI_SITE_URL"),
		"LikeName":        psychologist.LikeName,
		"PatientFullName": patient.FullName,
	})

	mail := &mails_models.TransientMailMessage{
		ID:          mailID,
		FromAddress: "relacionamento@psi.com.br",
		FromName:    "Relacionamento PSI",
		To:          psyUser.Email,
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
