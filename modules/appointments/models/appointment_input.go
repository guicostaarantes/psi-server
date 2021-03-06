package models

// EditAppointmentByPatientInput is the schema for information needed to edit an appointment by the patient
type EditAppointmentByPatientInput struct {
	Start  int64  `json:"start" bson:"start"`
	Reason string `json:"reason" bson:"reason"`
}

// EditAppointmentByPsychologistInput is the schema for information needed to edit an appointment by the psychologist
type EditAppointmentByPsychologistInput struct {
	Start  int64  `json:"start" bson:"start"`
	End    int64  `json:"end" bson:"end"`
	Price  int64  `json:"price" bson:"price"`
	Reason string `json:"reason" bson:"reason"`
}
