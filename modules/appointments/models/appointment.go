package appointments_models

import (
	"time"

	"gorm.io/gorm"
)

// AppointmentStatus represents the status of an appointment
type AppointmentStatus string

const (
	// Created means that the appointment was created by the jobrunner user
	Created AppointmentStatus = "CREATED"
	// ConfirmedByPatient means only that the patient has confirmed the appointment
	ConfirmedByPatient AppointmentStatus = "CONFIRMED_BY_PATIENT"
	// ConfirmedByPsychologist means only that the psychologist has confirmed the appointment
	ConfirmedByPsychologist AppointmentStatus = "CONFIRMED_BY_PSYCHOLOGIST"
	// ConfirmedByBoth means that both the patient and the psychologist have confirmed the appointment
	ConfirmedByBoth AppointmentStatus = "CONFIRMED_BY_BOTH"
	// EditedByPatient means that the patient has suggested a new time for the appointment and the psychologist needs to confirm
	EditedByPatient AppointmentStatus = "EDITED_BY_PATIENT"
	// EditedByPsychologist means that the psychologist has suggested a new time for the appointment and the patient needs to confirm
	EditedByPsychologist AppointmentStatus = "EDITED_BY_PSYCHOLOGIST"
	// CanceledByPatient means that the appointment was cancelled by the patient
	CanceledByPatient AppointmentStatus = "CANCELED_BY_PATIENT"
	// CanceledByPsychologist means that the appointment was cancelled by the psychologist
	CanceledByPsychologist AppointmentStatus = "CANCELED_BY_PSYCHOLOGIST"
	// TreatmentInterruptedByPatient means that the whole treatment was interrupted by the patient
	TreatmentInterruptedByPatient AppointmentStatus = "TREATMENT_INTERRUPTED_BY_PATIENT"
	// TreatmentInterruptedByPsychologist means that the whole treatment was interrupted by the psychologist
	TreatmentInterruptedByPsychologist AppointmentStatus = "TREATMENT_INTERRUPTED_BY_PSYCHOLOGIST"
	// TreatmentFinalized means that the whole treatment was finalized
	TreatmentFinalized AppointmentStatus = "TREATMENT_FINALIZED"
)

// Appointment represents the mutual promise of psychologist and patient to meet at a specific time
type Appointment struct {
	ID             string            `json:"id" gorm:"primaryKey"`
	CreatedAt      time.Time         `json:"createdAt`
	UpdatedAt      time.Time         `json:"updatedAt`
	DeletedAt      gorm.DeletedAt    `gorm:"index"`
	TreatmentID    string            `json:"treatmentId"`
	PatientID      string            `json:"patientId"`
	PsychologistID string            `json:"psychologistId"`
	Start          time.Time         `json:"start"`
	End            time.Time         `json:"end"`
	PriceRangeName string            `json:"priceRangeName"`
	Status         AppointmentStatus `json:"status"`
	Reason         string            `json:"reason"`
	Link           string            `json:"link"`
}
