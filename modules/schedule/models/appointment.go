package models

// AppointmentStatus represents the status of an appointment
type AppointmentStatus string

const (
	// Proposed means that the patient has chosen a time for the appointment and is waiting for the psychologist's confirmation
	Proposed AppointmentStatus = "PROPOSED"
	// Denied means that the psychologist denied the proposed appointment
	Denied AppointmentStatus = "DENIED"
	// Confirmed means that the psychologist has confirmed the proposed appointment
	Confirmed AppointmentStatus = "CONFIRMED"
	// CanceledByPatient means that the appointment has once been confirmed but then was cancelled by the patient
	CanceledByPatient AppointmentStatus = "CANCELED_BY_PATIENT"
	// CanceledByPsychologist means that the appointment has once been confirmed but then was cancelled by the psychologist
	CanceledByPsychologist AppointmentStatus = "CANCELED_BY_PSYCHOLOGIST"
)

// Appointment represents the mutual promise of psychologist and patient to consult at a specific time
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
