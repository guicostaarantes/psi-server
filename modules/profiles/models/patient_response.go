package models

// PatientCharacteristicResponse is the schema for a characteristic of a patient and its possible values to be returned to the user
type PatientCharacteristicResponse struct {
	Name           string   `json:"name" bson:"name"`
	Many           bool     `json:"many" bson:"many"`
	PossibleValues []string `json:"possibleValues" bson:"possibleValues"`
}

// PatientCharacteristicChoiceResponse is the schema for a characteristic of a patient and its possible values to be returned to the user
type PatientCharacteristicChoiceResponse struct {
	Name           string   `json:"name" bson:"name"`
	Many           bool     `json:"many" bson:"many"`
	Values         []string `json:"values" bson:"values"`
	PossibleValues []string `json:"possibleValues" bson:"possibleValues"`
}
