package models

// AffinityScore represents in the result field how much likely it is for a treatment to be succesful between psychologist and patient, based on their characteristics and preferences
type AffinityScore struct {
	ScoreForPatient      int64 `json:"scoreForPatient"`
	ScoreForPsychologist int64 `json:"scoreForPsychologist"`
}

// Affinity is the representation in the database of a calculation of affinity between psychologist and patient
type Affinity struct {
	PatientID            string `json:"patientId"`
	PsychologistID       string `json:"psychologistId"`
	CreatedAt            int64  `json:"createdAt"`
	ScoreForPatient      int64  `json:"scoreForPatient"`
	ScoreForPsychologist int64  `json:"scoreForPsychologist"`
}
