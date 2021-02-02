package profiles_models

// Psychologist is the schema for the profile of a psychologist
type Psychologist struct {
	ID        string `json:"id" bson:"id"`
	UserID    string `json:"userId" bson:"userId"`
	BirthDate int64  `json:"birthDate" bson:"birthDate"`
	City      string `json:"city" bson:"city"`
}

// PsyCharacteristic is the schema for a characteristic of a psychologist and its possible values
type PsyCharacteristic struct {
	ID     string `json:"id" bson:"id"`
	Name   string `json:"name" bson:"name"`
	Many   bool   `json:"many" bson:"many"`
	Values string `json:"values" bson:"values"`
}

// PsyChoice is the schema for a choice of characteristic from a psychologist
type PsyChoice struct {
	PsychologistID      string `json:"psychologistId" bson:"psychologistId"`
	PsyCharacteristicID string `json:"psyCharacteristicId" bson:"psyCharacteristicId"`
	Value               string `json:"value" bson:"value"`
}

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

// CreatePsyCharacteristicInput is the schema for information needed to create a characteristic of a psychologist and its possible values
type CreatePsyCharacteristicInput struct {
	Name   string   `json:"name" bson:"name"`
	Many   bool     `json:"many" bson:"many"`
	Values []string `json:"values" bson:"values"`
}

// UpdatePsyCharacteristicInput is the schema for information needed to create a characteristic of a psychologist and its possible values
type UpdatePsyCharacteristicInput struct {
	Name   string   `json:"name" bson:"name"`
	Many   bool     `json:"many" bson:"many"`
	Values []string `json:"values" bson:"values"`
}

// AddPsyChoice is the schema for adding a choice of characteristic from a psychologist
type AddPsyChoice struct {
	PsychologistID      string `json:"psychologistId" bson:"psychologistId"`
	PsyCharacteristicID string `json:"psyCharacteristicId" bson:"psyCharacteristicId"`
	Value               string `json:"value" bson:"value"`
}

// RemovePsyChoice is the schema for removing a choice of characteristic from a psychologist
type RemovePsyChoice struct {
	PsychologistID      string `json:"psychologistId" bson:"psychologistId"`
	PsyCharacteristicID string `json:"psyCharacteristicId" bson:"psyCharacteristicId"`
	Value               string `json:"value" bson:"value"`
}
