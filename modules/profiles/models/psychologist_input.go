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
