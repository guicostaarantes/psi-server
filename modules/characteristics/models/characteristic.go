package models

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
	Name           string               `json:"name" bson:"name"`
	Type           CharacteristicType   `json:"type" bson:"type"`
	Target         CharacteristicTarget `json:"target" bson:"target"`
	PossibleValues string               `json:"possibleValues" bson:"possibleValues"`
}

// CharacteristicChoice is the schema for a choice of characteristics made by a profile
type CharacteristicChoice struct {
	ProfileID          string `json:"profileId" bson:"profileId"`
	CharacteristicName string `json:"characteristicName" bson:"characteristicName"`
	SelectedValue      string `json:"selectedValue" bson:"selectedValue"`
}

// Preference is the schema for the fact that a patient prefers working with a certain kind of psychologist, and vice-versa
type Preference struct {
	ProfileID          string `json:"profileId" bson:"profileId"`
	CharacteristicName string `json:"characteristicName" bson:"characteristicName"`
	SelectedValue      string `json:"selectedValue" bson:"selectedValue"`
	Weight             int64  `json:"weight" bson:"weight"`
}
