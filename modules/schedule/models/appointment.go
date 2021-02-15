package models

// AppointmentStatus represents the status of an appointment
type AppointmentStatus string

const (
	// Proposed means that the patient has chosen a time for the appointment and is waiting for the psychologist's confirmation
	Proposed AppointmentStatus = "PROPOSED"
	// AwaitingRescheduling means that the psychologist denied the proposed appointment and is waiting for a new proposal
	AwaitingRescheduling AppointmentStatus = "AWAITING_RESCHEDULING"
	// Confirmed means that the psychologist has confirmed the proposed appointment
	Confirmed AppointmentStatus = "CONFIRMED"
	// Canceled means that the appointment has once been confirmed but then was cancelled by one of the parts
	Canceled AppointmentStatus = "CANCELED"
)

// Appointment represents the mutual promise of psychologist and patient to consult at a specific time
type Appointment struct {
	ID          string            `json:"id" bson:"id"`
	TreatmentID string            `json:"treatmentId" bson:"treatmentId"`
	Start       int64             `json:"start" bson:"start"`
	End         int64             `json:"end" bson:"end"`
	Price       int64             `json:"price" bson:"price"`
	Status      AppointmentStatus `json:"status" bson:"status"`
}
