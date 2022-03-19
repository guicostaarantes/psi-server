package treatments_services

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"text/template"
	"time"

	appointments_models "github.com/guicostaarantes/psi-server/modules/appointments/models"
	mails_models "github.com/guicostaarantes/psi-server/modules/mails/models"
	profiles_models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	treatments_models "github.com/guicostaarantes/psi-server/modules/treatments/models"
	treatments_templates "github.com/guicostaarantes/psi-server/modules/treatments/templates"
	users_models "github.com/guicostaarantes/psi-server/modules/users/models"
	"github.com/guicostaarantes/psi-server/utils/identifier"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// InterruptTreatmentByPsychologistService is a service that interrupts a treatment, changing its status to interrupted by psychologist
type InterruptTreatmentByPsychologistService struct {
	IdentifierUtil identifier.IIdentifierUtil
	OrmUtil        orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s InterruptTreatmentByPsychologistService) Execute(id string, psychologistID string, reason string) error {

	treatment := treatments_models.Treatment{}
	psychologist := profiles_models.Psychologist{}
	patient := profiles_models.Patient{}
	patientUser := users_models.User{}

	result := s.OrmUtil.Db().Where("id = ? AND psychologist_id = ?", id, psychologistID).Limit(1).Find(&treatment)
	if result.Error != nil {
		return result.Error
	}

	if treatment.ID == "" {
		return errors.New("resource not found")
	}

	result = s.OrmUtil.Db().Where("id = ?", psychologistID).Limit(1).Find(&psychologist)
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

	if treatment.Status != treatments_models.Active {
		return fmt.Errorf("treatments can only be interrupted if their current status is ACTIVE. current status is %s", string(treatment.Status))
	}

	appointments := []*appointments_models.Appointment{}

	result = s.OrmUtil.Db().Where("treatment_id = ? AND start > ?", id, time.Now()).Find(&appointments)
	if result.Error != nil {
		return result.Error
	}

	for _, appointment := range appointments {
		if appointment.Status != appointments_models.CanceledByPatient {
			appointment.Status = appointments_models.TreatmentInterruptedByPsychologist
			appointment.Reason = reason

			result = s.OrmUtil.Db().Save(&appointment)
			if result.Error != nil {
				return result.Error
			}
		}
	}

	now := time.Now()
	treatment.EndDate = &now
	treatment.Status = treatments_models.InterruptedByPsychologist
	treatment.Reason = reason

	_, mailID, mailIDErr := s.IdentifierUtil.GenerateIdentifier()
	if mailIDErr != nil {
		return mailIDErr
	}

	templ, templErr := template.New("TreatmentInterruptedByPsychologistEmail").Parse(treatments_templates.TreatmentInterruptedByPsychologistEmailTemplate)
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
		Subject:     "Tratamento interrompido no PSI",
		Html:        buff.String(),
		Processed:   false,
	}

	result = s.OrmUtil.Db().Create(&mail)
	if result.Error != nil {
		return result.Error
	}

	result = s.OrmUtil.Db().Save(&treatment)
	if result.Error != nil {
		return result.Error
	}

	return nil

}
