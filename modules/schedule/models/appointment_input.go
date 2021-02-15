package models

// ProposeAppointmentInput is the schema for information needed to propose an appointment
type ProposeAppointmentInput struct {
	TreatmentID string `json:"treatmentId" bson:"treatmentId"`
	Start       int64  `json:"start" bson:"start"`
}
