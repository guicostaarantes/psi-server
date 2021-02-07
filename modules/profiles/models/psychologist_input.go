package models

// CreatePsychologistInput is the schema for information needed to create a psychologist
type CreatePsychologistInput struct {
	UserID    string `json:"userId" bson:"userId"`
	BirthDate int64  `json:"birthDate" bson:"birthDate"`
	City      string `json:"city" bson:"city"`
}

// UpdatePsychologistInput is the schema for information needed to update a psychologist
type UpdatePsychologistInput struct {
	BirthDate int64  `json:"birthDate" bson:"birthDate"`
	City      string `json:"city" bson:"city"`
}

// SetPsychologistCharacteristicInput is the schema for information needed to create a characteristic of a psychologist and its possible values
type SetPsychologistCharacteristicInput struct {
	Name           string   `json:"name" bson:"name"`
	Many           bool     `json:"many" bson:"many"`
	PossibleValues []string `json:"possibleValues" bson:"possibleValues"`
}

// SetPsychologistCharacteristicChoiceInput is the schema for information needed to assign a characteristic to a psychologist profile
type SetPsychologistCharacteristicChoiceInput struct {
	PsychologistID     string   `json:"psychologistId" bson:"psychologistId"`
	CharacteristicName string   `json:"characteristicName" bson:"characteristicName"`
	Values             []string `json:"values" bson:"values"`
}
