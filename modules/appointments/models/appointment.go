package models

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
)

// Appointment represents the mutual promise of psychologist and patient to meet at a specific time
type Appointment struct {
	ID             string            `json:"id" bson:"id"`
	TreatmentID    string            `json:"treatmentId" bson:"treatmentId"`
	PatientID      string            `json:"patientId" bson:"patientId"`
	PsychologistID string            `json:"psychologistId" bson:"psychologistId"`
	Start          int64             `json:"start" bson:"start"`
	End            int64             `json:"end" bson:"end"`
	Price          int64             `json:"price" bson:"price"`
	Status         AppointmentStatus `json:"status" bson:"status"`
	Reason         string            `json:"reason" bson:"reason"`
	Link           string            `json:"link" bson:"link"`
}
