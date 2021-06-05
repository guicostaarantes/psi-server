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

// Treatment represents the intention from a psychologist to treat a patient, defining the session weeklyStart, duration and price of sessions
// All treatments have weekly sessions, and WeeklyStart is the number of seconds from Sundays at midnight where the sessions happen
type Treatment struct {
	ID             string          `json:"id" bson:"id"`
	PsychologistID string          `json:"psychologistId" bson:"psychologistId"`
	PatientID      string          `json:"patientId" bson:"patientId"`
	WeeklyStart    int64           `json:"weeklyStart" bson:"weeklyStart"`
	Duration       int64           `json:"duration" bson:"duration"`
	Price          int64           `json:"price" bson:"price"`
	Status         TreatmentStatus `json:"status" bson:"status"`
	StartDate      int64           `json:"startDate" bson:"startDate"`
	EndDate        int64           `json:"endDate" bson:"endDate"`
	Reason         string          `json:"reason" bson:"reason"`
}
