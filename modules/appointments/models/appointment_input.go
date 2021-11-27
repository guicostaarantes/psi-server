package appointments_models

import "time"

// EditAppointmentByPatientInput is the schema for information needed to edit an appointment by the patient
type EditAppointmentByPatientInput struct {
	Start  time.Time `json:"start"`
	Reason string    `json:"reason"`
}

// EditAppointmentByPsychologistInput is the schema for information needed to edit an appointment by the psychologist
type EditAppointmentByPsychologistInput struct {
	Start          time.Time `json:"start"`
	End            time.Time `json:"end"`
	PriceRangeName string    `json:"priceRangeName"`
	Reason         string    `json:"reason"`
}
