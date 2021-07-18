package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// AssignTreatmentService is a service that assigns a patient to a treatment and changes its status to active
type AssignTreatmentService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s AssignTreatmentService) Execute(id string, patientID string) error {

	treatment := models.Treatment{}

	patientInOtherTreatment := models.Treatment{}

	findErr := s.DatabaseUtil.FindOne("treatments", map[string]interface{}{"patientId": patientID, "status": string(models.Active)}, &patientInOtherTreatment)
	if findErr != nil {
		return findErr
	}

	if patientInOtherTreatment.ID != "" {
		return errors.New("patient is already in an active treatment")
	}

	findErr = s.DatabaseUtil.FindOne("treatments", map[string]interface{}{"id": id}, &treatment)
	if findErr != nil {
		return findErr
	}

	if treatment.ID == "" {
		return errors.New("resource not found")
	}

	if treatment.Status != models.Pending {
		return fmt.Errorf("treatments can only be assigned if their current status is PENDING. current status is %s", string(treatment.Status))
	}

	treatment.PatientID = patientID
	treatment.StartDate = time.Now().Unix()
	treatment.Status = models.Active

	writeErr := s.DatabaseUtil.UpdateOne("treatments", map[string]interface{}{"id": id}, treatment)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
