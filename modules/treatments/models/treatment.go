package models

// TreatmentStatus represents the status of a treatment
type TreatmentStatus string

const (
	// Pending means that the treatment has been created but no patient has yet occupied it
	Pending TreatmentStatus = "PENDING"
	// Active means that the treatment has a patient and the treatment is occuring
	Active TreatmentStatus = "ACTIVE"
	// Finalized means that the treatment had a patient and the treatment was finished succesfully
	Finalized TreatmentStatus = "FINALIZED"
	// InterruptedByPsychologist means that the treatment had a patient and the treatment was interrupted by the psychologist
	InterruptedByPsychologist TreatmentStatus = "INTERRUPTED_BY_PSYCHOLOGIST"
	// InterruptedByPatient means that the treatment had a patient and the treatment was interrupted by the patient
	InterruptedByPatient TreatmentStatus = "INTERRUPTED_BY_PATIENT"
)

// Treatment represents the intention from a psychologist to treat a patient, defining the sessions' duration, price, interval and phase.
// The next session of a specific treatment will be scheduled to the UNIX timestamp T, where T = (ScheduleIntervalSeconds * Frequency * N) + Phase, and N is the smallest natural number that makes T superior to the current timestamp.
type Treatment struct {
	ID             string          `json:"id" bson:"id"`
	PsychologistID string          `json:"psychologistId" bson:"psychologistId"`
	PatientID      string          `json:"patientId" bson:"patientId"`
	Frequency      int64           `json:"frequency" bson:"frequency"`
	Phase          int64           `json:"phase" bson:"phase"`
	Duration       int64           `json:"duration" bson:"duration"`
	PriceRange     string          `json:"priceRange" bson:"priceRange"`
	Status         TreatmentStatus `json:"status" bson:"status"`
	StartDate      int64           `json:"startDate" bson:"startDate"`
	EndDate        int64           `json:"endDate" bson:"endDate"`
	Reason         string          `json:"reason" bson:"reason"`
}
