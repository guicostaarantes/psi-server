package treatments_services

import (
	"bytes"
	"errors"
	"os"
	"text/template"

	mails_models "github.com/guicostaarantes/psi-server/modules/mails/models"
	profiles_models "github.com/guicostaarantes/psi-server/modules/profiles/models"
	treatments_models "github.com/guicostaarantes/psi-server/modules/treatments/models"
	treatments_templates "github.com/guicostaarantes/psi-server/modules/treatments/templates"
	users_models "github.com/guicostaarantes/psi-server/modules/users/models"
	"github.com/guicostaarantes/psi-server/utils/identifier"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// UpdateTreatmentService is a service that changes data from a treatment
type UpdateTreatmentService struct {
	IdentifierUtil                 identifier.IIdentifierUtil
	OrmUtil                        orm.IOrmUtil
	CheckTreatmentCollisionService *CheckTreatmentCollisionService
}

// Execute is the method that runs the business logic of the service
func (s UpdateTreatmentService) Execute(id string, psychologistID string, input treatments_models.UpdateTreatmentInput) error {

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

	if treatment.PatientID != "" {

		result = s.OrmUtil.Db().Where("id = ?", treatment.PatientID).Limit(1).Find(&patient)
		if result.Error != nil {
			return result.Error
		}

		result = s.OrmUtil.Db().Where("id = ?", patient.UserID).Limit(1).Find(&patientUser)
		if result.Error != nil {
			return result.Error
		}

	}

	if input.PriceRangeName != "" && treatment.Status == treatments_models.Pending {
		return errors.New("pending treatments are not allowed to have a price range")
	}

	if input.PriceRangeName == "" && treatment.Status != treatments_models.Pending {
		return errors.New("non-pending treatments must have a price range")
	}

	checkErr := s.CheckTreatmentCollisionService.Execute(psychologistID, input.Frequency, input.Phase, input.Duration, id)
	if checkErr != nil {
		return checkErr
	}

	treatment.Frequency = input.Frequency
	treatment.Phase = input.Phase
	treatment.Duration = input.Duration
	treatment.PriceRangeName = input.PriceRangeName

	if treatment.PatientID != "" {

		_, mailID, mailIDErr := s.IdentifierUtil.GenerateIdentifier()
		if mailIDErr != nil {
			return mailIDErr
		}

		templ, templErr := template.New("TreatmentModifiedEmail").Parse(treatments_templates.TreatmentModifiedEmailTemplate)
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
			Subject:     "Tratamento modificado no PSI",
			Html:        buff.String(),
			Processed:   false,
		}

		result = s.OrmUtil.Db().Create(&mail)
		if result.Error != nil {
			return result.Error
		}
	}

	result = s.OrmUtil.Db().Save(&treatment)
	if result.Error != nil {
		return result.Error
	}

	return nil

}
