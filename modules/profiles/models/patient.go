package models

// Patient is the schema for the profile of a patient
type Patient struct {
	ID        string `json:"id" bson:"id"`
	UserID    string `json:"userId" bson:"userId"`
	BirthDate int64  `json:"birthDate" bson:"birthDate"`
	City      string `json:"city" bson:"city"`
}

// PatientCharacteristic is the schema for a characteristic of a patient and its possible values
type PatientCharacteristic struct {
	Name           string `json:"name" bson:"name"`
	Many           bool   `json:"many" bson:"many"`
	PossibleValues string `json:"possibleValues" bson:"possibleValues"`
}

// PatientCharacteristicChoice is the schema for a choice of characteristic from a patient
type PatientCharacteristicChoice struct {
	PatientID          string `json:"patientId" bson:"patientId"`
	CharacteristicName string `json:"characteristicName" bson:"characteristicName"`
	Value              string `json:"value" bson:"value"`
}

// PatientPreference is the schema for the fact that a patient prefers working with a certain kind of psychologist
type PatientPreference struct {
	PatientID          string `json:"patientId" bson:"patientId"`
	CharacteristicName string `json:"characteristicName" bson:"characteristicName"`
	Value              string `json:"value" bson:"value"`
	Weight             int64  `json:"weight" bson:"weight"`
}
