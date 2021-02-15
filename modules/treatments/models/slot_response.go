package models

// GetPsychologistTreatmentsResponse is the schema for information needed to be sent to the psychologist about their treatments
type GetPsychologistTreatmentsResponse struct {
	ID        string          `json:"id" bson:"id"`
	PatientID string          `json:"patientId" bson:"patientId"`
	Duration  int64           `json:"duration" bson:"duration"`
	Price     int64           `json:"price" bson:"price"`
	Interval  int64           `json:"interval" bson:"interval"`
	Status    TreatmentStatus `json:"status" bson:"status"`
}

// GetPatientTreatmentsResponse is the schema for information needed to be sent to the patient about their treatments
type GetPatientTreatmentsResponse struct {
	ID             string          `json:"id" bson:"id"`
	PsychologistID string          `json:"psychologistId" bson:"psychologistId"`
	Duration       int64           `json:"duration" bson:"duration"`
	Price          int64           `json:"price" bson:"price"`
	Interval       int64           `json:"interval" bson:"interval"`
	Status         TreatmentStatus `json:"status" bson:"status"`
}
