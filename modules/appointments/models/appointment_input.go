package models

// EditAppointmentByPatientInput is the schema for information needed to edit an appointment by the patient
type EditAppointmentByPatientInput struct {
	Start  int64  `json:"start" bson:"start"`
	Reason string `json:"reason" bson:"reason"`
}

// EditAppointmentByPsychologistInput is the schema for information needed to edit an appointment by the psychologist
type EditAppointmentByPsychologistInput struct {
	Start          int64  `json:"start" bson:"start"`
	End            int64  `json:"end" bson:"end"`
	PriceRangeName string `json:"priceRangeName" bson:"priceRangeName"`
	Reason         string `json:"reason" bson:"reason"`
}
