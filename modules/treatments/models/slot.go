package models

// SlotStatus represents the status of a slot
type SlotStatus string

const (
	// Pending means that the slot has been created but no patient has yet occupied it
	Pending SlotStatus = "PENDING"
	// Active means that the slot has a patient and the treatment is occuring
	Active SlotStatus = "ACTIVE"
	// Finalized means that the slot had a patient and the treatment was finished succesfully
	Finalized SlotStatus = "FINALIZED"
	// InterruptedByPsychologist means that the slot had a patient and the treatment was interrupted by the psychologist
	InterruptedByPsychologist SlotStatus = "INTERRUPTED_BY_PSYCHOLOGIST"
	// InterruptedByPatient means that the slot had a patient and the treatment was interrupted by the patient
	InterruptedByPatient SlotStatus = "INTERRUPTED_BY_PATIENT"
)

// Slot represents the intention from a psychologist to treat a patient, defining the session duration, price and interval between sessions
type Slot struct {
	ID             string     `json:"id" bson:"id"`
	PsychologistID string     `json:"psychologistId" bson:"psychologistId"`
	PatientID      string     `json:"patientId" bson:"patientId"`
	Duration       int64      `json:"duration" bson:"duration"`
	Price          int64      `json:"price" bson:"price"`
	Interval       int64      `json:"interval" bson:"interval"`
	Status         SlotStatus `json:"status" bson:"status"`
	StartDate      int64      `json:"startDate" bson:"startDate"`
	EndDate        int64      `json:"endDate" bson:"endDate"`
	Reason         string     `json:"reason" bson:"reason"`
}
