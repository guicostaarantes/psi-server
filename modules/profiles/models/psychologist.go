package models

// Psychologist is the schema for the profile of a psychologist
type Psychologist struct {
	ID        string `json:"id" bson:"id"`
	UserID    string `json:"userId" bson:"userId"`
	BirthDate int64  `json:"birthDate" bson:"birthDate"`
	City      string `json:"city" bson:"city"`
}

// PsychologistCharacteristic is the schema for a characteristic of a psychologist and its possible values
type PsychologistCharacteristic struct {
	Name           string `json:"name" bson:"name"`
	Many           bool   `json:"many" bson:"many"`
	PossibleValues string `json:"possibleValues" bson:"possibleValues"`
}

// PsychologistCharacteristicChoice is the schema for a choice of characteristic from a psychologist
type PsychologistCharacteristicChoice struct {
	PsychologistID     string `json:"psychologistId" bson:"psychologistId"`
	CharacteristicName string `json:"characteristicName" bson:"characteristicName"`
	Value              string `json:"value" bson:"value"`
}

// PsychologistPreference is the schema for the fact that a psychologist prefers working with a certain kind of patient
type PsychologistPreference struct {
	PsychologistID     string `json:"psychologistId" bson:"psychologistId"`
	CharacteristicName string `json:"characteristicName" bson:"characteristicName"`
	Value              string `json:"value" bson:"value"`
	Weight             int64  `json:"weight" bson:"weight"`
}
