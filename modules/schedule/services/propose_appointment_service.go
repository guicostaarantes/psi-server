package services

import (
	"context"
	"errors"

	"github.com/guicostaarantes/psi-server/modules/schedule/models"
	treatments_models "github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/identifier"
)

// ProposeAppointmentService is a service that the patient will use to ask the psychologist for an appointment
type ProposeAppointmentService struct {
	DatabaseUtil   database.IDatabaseUtil
	IdentifierUtil identifier.IIdentifierUtil
}

// Execute is the method that runs the business logic of the service
func (s ProposeAppointmentService) Execute(patientID string, input models.ProposeAppointmentInput) error {

	treatment := treatments_models.Treatment{}

	findErr := s.DatabaseUtil.FindOne("psi_db", "treatments", map[string]interface{}{"id": input.TreatmentID, "patientId": patientID, "status": string(treatments_models.Active)}, &treatment)
	if findErr != nil {
		return findErr
	}

	if treatment.ID == "" {
		return errors.New("resource not found")
	}

	otherAppointment := models.Appointment{}

	findErr = s.DatabaseUtil.FindOne("psi_db", "appointments", map[string]interface{}{"patientId": patientID, "status": string(models.Proposed)}, &otherAppointment)
	if findErr != nil {
		return findErr
	}

	if otherAppointment.ID != "" {
		return errors.New("patient already has an appointment with status PROPOSED")
	}

	end := input.Start + treatment.Duration

	isAvailable := false

	avalCursor, findAvalErr := s.DatabaseUtil.FindMany("psi_db", "availabilities", map[string]interface{}{"psychologistId": treatment.PsychologistID})
	if findAvalErr != nil {
		return findAvalErr
	}

	defer avalCursor.Close(context.Background())

	for avalCursor.Next(context.Background()) {

		aval := models.Availability{}

		decodeErr := avalCursor.Decode(&aval)
		if decodeErr != nil {
			return decodeErr
		}

		if aval.Start <= input.Start && aval.End >= end {
			isAvailable = true
			break
		}

	}

	if !isAvailable {
		return errors.New("the psychologist is not available during the requested time slot")
	}

	isClash := false

	clashCursor, findClashErr := s.DatabaseUtil.FindMany("psi_db", "appointments", map[string]interface{}{"psychologistId": treatment.PsychologistID, "status": string(models.Confirmed)})
	if findClashErr != nil {
		return findClashErr
	}

	defer clashCursor.Close(context.Background())

	for clashCursor.Next(context.Background()) {

		appointment := models.Appointment{}

		decodeErr := clashCursor.Decode(&appointment)
		if decodeErr != nil {
			return decodeErr
		}

		if appointment.End > input.Start && appointment.Start < end {
			isClash = true
			break
		}

	}

	if isClash {
		return errors.New("the psychologist is not available during the requested time slot")
	}

	_, appoID, appoIDErr := s.IdentifierUtil.GenerateIdentifier()
	if appoIDErr != nil {
		return appoIDErr
	}

	newAppointment := models.Appointment{
		ID:             appoID,
		TreatmentID:    input.TreatmentID,
		PatientID:      treatment.PatientID,
		PsychologistID: treatment.PsychologistID,
		Start:          input.Start,
		End:            end,
		Price:          treatment.Price,
		Status:         models.Proposed,
	}

	writeErr := s.DatabaseUtil.InsertOne("psi_db", "appointments", newAppointment)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
