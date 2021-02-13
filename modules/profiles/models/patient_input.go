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
