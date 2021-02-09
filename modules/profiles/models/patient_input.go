package models

// CreatePatientInput is the schema for information needed to create a patient
type CreatePatientInput struct {
	UserID    string `json:"userId" bson:"userId"`
	BirthDate int64  `json:"birthDate" bson:"birthDate"`
	City      string `json:"city" bson:"city"`
}

// UpdatePatientInput is the schema for information needed to update a patient
type UpdatePatientInput struct {
	BirthDate int64  `json:"birthDate" bson:"birthDate"`
	City      string `json:"city" bson:"city"`
}

// SetPatientCharacteristicInput is the schema for information needed to create a characteristic of a patient and its possible values
type SetPatientCharacteristicInput struct {
	Name           string   `json:"name" bson:"name"`
	Many           bool     `json:"many" bson:"many"`
	PossibleValues []string `json:"possibleValues" bson:"possibleValues"`
}

// SetPatientCharacteristicChoiceInput is the schema for information needed to assign a characteristic to a patient profile
type SetPatientCharacteristicChoiceInput struct {
	CharacteristicName string   `json:"characteristicName" bson:"characteristicName"`
	Values             []string `json:"values" bson:"values"`
}

// SetPatientPreferenceInput is the schema for information needed to set the preferences of a patient
type SetPatientPreferenceInput struct {
	CharacteristicName string `json:"characteristicName" bson:"characteristicName"`
	Value              string `json:"value" bson:"value"`
	Weight             int64  `json:"weight" bson:"weight"`
}
