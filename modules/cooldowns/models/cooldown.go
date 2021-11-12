package models

import (
	"time"
)

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
	// TreatmentInterrupted means that the user has interrupted a treatment
	TreatmentInterrupted CooldownType = "TREATMENT_INTERRUPTED"
	// TopAffinitiesSet means that the user set their top affinities for a psychologist
	TopAffinitiesSet CooldownType = "TOP_AFFINITIES_SET"
)

// Cooldown holds information about the usage of the system
type Cooldown struct {
	ID           string              `json:"id" gorm:"primaryKey"`
	CreatedAt    time.Time           `json:"createdAt"`
	UpdatedAt    time.Time           `json:"updatedAt"`
	ProfileID    string              `json:"profileId"`
	ProfileType  CooldownProfileType `json:"profileType"`
	CooldownType CooldownType        `json:"cooldownType"`
	ValidUntil   int64               `json:"validUntil"`
}
