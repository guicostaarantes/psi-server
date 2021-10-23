package models

// TermProfileType represents the possible profiles responsible for an agreement
type TermProfileType string

const (
	// PatientTarget means that the agreement is related to a patient
	Patient TermProfileType = "PATIENT"
	// PsychologistTarget means that the agreement is related to a psychologist
	Psychologist TermProfileType = "PSYCHOLOGIST"
)

// Term is a group of rules and instructions that users should agree in order to use the platform
type Term struct {
	Name        string          `json:"name" bson:"name"`
	Version     int64           `json:"version" bson:"version"`
	ProfileType TermProfileType `json:"profileType" bson:"profileType"`
	Active      bool            `json:"active" bson:"active"`
}

// Legible content of the term should be available in translations collection under key {ProfileType (psy | pat)}-term:{Name}:{Version}
