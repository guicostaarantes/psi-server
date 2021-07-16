package models

// CooldownProfileType represents the possible profiles responsible for a cooldown
type CooldownProfileType string

const (
	// PatientTarget means that the cooldown is related to a patient
	Patient CooldownProfileType = "PATIENT"
	// PsychologistTarget means that the cooldown is related to a psychologist
	Psychologist CooldownProfileType = "PSYCHOLOGIST"
)

// CooldownType represents the possible types of a cooldown
type CooldownType string

const (
	// TopAffinitiesSet means that the user set their top affinities for a psychologist
	TopAffinitiesSet CooldownType = "TOP_AFFINITIES_SET"
)

// Cooldown holds information about the usage of the system
type Cooldown struct {
	ID           string              `json:"id" bson:"id"`
	ProfileID    string              `json:"profileId" bson:"profileId"`
	ProfileType  CooldownProfileType `json:"profileType" bson:"profileType"`
	CooldownType CooldownType        `json:"cooldownType" bson:"cooldownType"`
	CreatedAt    int64               `json:"createdAt" bson:"createdAt"`
	ValidUntil   int64               `json:"validUntil" bson:"validUntil"`
}
