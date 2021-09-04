package models

// GetPsychologistTreatmentsResponse is the schema for information needed to be sent to the psychologist about their treatments
type GetPsychologistTreatmentsResponse struct {
	ID             string          `json:"id" bson:"id"`
	PatientID      string          `json:"patientId" bson:"patientId"`
	Frequency      int64           `json:"frequency" bson:"frequency"`
	Phase          int64           `json:"phase" bson:"phase"`
	Duration       int64           `json:"duration" bson:"duration"`
	PriceRangeName string          `json:"priceRangeName" bson:"priceRangeName"`
	Status         TreatmentStatus `json:"status" bson:"status"`
}

// GetPatientTreatmentsResponse is the schema for information needed to be sent to the patient about their treatments
type GetPatientTreatmentsResponse struct {
	ID             string          `json:"id" bson:"id"`
	PsychologistID string          `json:"psychologistId" bson:"psychologistId"`
	Frequency      int64           `json:"frequency" bson:"frequency"`
	Phase          int64           `json:"phase" bson:"phase"`
	Duration       int64           `json:"duration" bson:"duration"`
	PriceRangeName string          `json:"priceRangeName" bson:"priceRangeName"`
	Status         TreatmentStatus `json:"status" bson:"status"`
}
