package appointments_models

// EditAppointmentByPatientInput is the schema for information needed to edit an appointment by the patient
type EditAppointmentByPatientInput struct {
	Start  int64  `json:"start"`
	Reason string `json:"reason"`
}

// EditAppointmentByPsychologistInput is the schema for information needed to edit an appointment by the psychologist
type EditAppointmentByPsychologistInput struct {
	Start          int64  `json:"start"`
	End            int64  `json:"end"`
	PriceRangeName string `json:"priceRangeName"`
	Reason         string `json:"reason"`
}
