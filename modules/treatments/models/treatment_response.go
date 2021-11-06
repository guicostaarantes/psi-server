package models

// GetPsychologistTreatmentsResponse is the schema for information needed to be sent to the psychologist about their treatments
type GetPsychologistTreatmentsResponse struct {
	ID             string          `json:"id"`
	PatientID      string          `json:"patientId"`
	Frequency      int64           `json:"frequency"`
	Phase          int64           `json:"phase"`
	Duration       int64           `json:"duration"`
	PriceRangeName string          `json:"priceRangeName"`
	Status         TreatmentStatus `json:"status"`
}

// GetPatientTreatmentsResponse is the schema for information needed to be sent to the patient about their treatments
type GetPatientTreatmentsResponse struct {
	ID             string          `json:"id"`
	PsychologistID string          `json:"psychologistId"`
	Frequency      int64           `json:"frequency"`
	Phase          int64           `json:"phase"`
	Duration       int64           `json:"duration"`
	PriceRangeName string          `json:"priceRangeName"`
	Status         TreatmentStatus `json:"status"`
}
