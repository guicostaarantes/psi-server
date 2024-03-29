package characteristics_models

import (
	"time"

	"gorm.io/gorm"
)

// CharacteristicType represents the possible inputs and choices of a characteristic
type CharacteristicType string

const (
	// Boolean is a type for a characteristic that can only be either true or false
	Boolean CharacteristicType = "BOOLEAN"
	// Single is a type for a characteristic that has multiple options but only one choice
	Single CharacteristicType = "SINGLE"
	// Multiple is a type for a characteristic that has multiple options and may have zero, one or multiple choices
	Multiple CharacteristicType = "MULTIPLE"
)

// CharacteristicTarget represents the possible receivers of a characteristic
type CharacteristicTarget string

const (
	// PatientTarget means that the characteristic is related to a patient
	PatientTarget CharacteristicTarget = "PATIENT"
	// PsychologistTarget means that the characteristic is related to a psychologist
	PsychologistTarget CharacteristicTarget = "PSYCHOLOGIST"
)

// Characteristic is the schema for a characteristic and its possible values
type Characteristic struct {
	ID             string               `json:"id" gorm:"primaryKey"`
	CreatedAt      time.Time            `json:"createdAt`
	UpdatedAt      time.Time            `json:"updatedAt`
	DeletedAt      gorm.DeletedAt       `gorm:"index"`
	Name           string               `json:"name"`
	Type           CharacteristicType   `json:"type"`
	Target         CharacteristicTarget `json:"target"`
	PossibleValues string               `json:"possibleValues"`
}

// CharacteristicChoice is the schema for a choice of characteristics made by a profile
type CharacteristicChoice struct {
	ID                 string               `json:"id" gorm:"primaryKey"`
	ProfileID          string               `json:"profileId" gorm:"index"`
	Target             CharacteristicTarget `json:"target"`
	CharacteristicName string               `json:"characteristicName"`
	SelectedValue      string               `json:"selectedValue"`
}

// Preference is the schema for the fact that a patient prefers working with a certain kind of psychologist, and vice-versa
type Preference struct {
	ProfileID          string               `json:"profileId" gorm:"index"`
	Target             CharacteristicTarget `json:"target" gorm:"index"`
	CharacteristicName string               `json:"characteristicName"`
	SelectedValue      string               `json:"selectedValue"`
	Weight             int64                `json:"weight"`
}
