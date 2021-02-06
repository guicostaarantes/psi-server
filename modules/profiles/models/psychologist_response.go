package models

// PsychologistCharacteristicResponse is the schema for a characteristic of a psychologist and its possible values to be returned to the user
type PsychologistCharacteristicResponse struct {
	ID             string   `json:"id" bson:"id"`
	Name           string   `json:"name" bson:"name"`
	Many           bool     `json:"many" bson:"many"`
	PossibleValues []string `json:"possibleValues" bson:"possibleValues"`
}

// PsychologistCharacteristicChoiceResponse is the schema for a characteristic of a psychologist and its possible values to be returned to the user
type PsychologistCharacteristicChoiceResponse struct {
	ID             string   `json:"id" bson:"id"`
	Name           string   `json:"name" bson:"name"`
	Many           bool     `json:"many" bson:"many"`
	Values         []string `json:"values" bson:"values"`
	PossibleValues []string `json:"possibleValues" bson:"possibleValues"`
}
