package models

// Patient is the schema for the profile of a patient
type Patient struct {
	ID        string `json:"id" bson:"id"`
	UserID    string `json:"userId" bson:"userId"`
	BirthDate int64  `json:"birthDate" bson:"birthDate"`
	City      string `json:"city" bson:"city"`
}
