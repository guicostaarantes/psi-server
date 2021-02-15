package models

// CreatePatientInput is the schema for information needed to create a patient
type CreatePatientInput struct {
	UserID    string `json:"userId" bson:"userId"`
	FullName  string `json:"fullName" bson:"fullName"`
	LikeName  string `json:"likeName" bson:"likeName"`
	BirthDate int64  `json:"birthDate" bson:"birthDate"`
	City      string `json:"city" bson:"city"`
}

// UpdatePatientInput is the schema for information needed to update a patient
type UpdatePatientInput struct {
	FullName  string `json:"fullName" bson:"fullName"`
	LikeName  string `json:"likeName" bson:"likeName"`
	BirthDate int64  `json:"birthDate" bson:"birthDate"`
	City      string `json:"city" bson:"city"`
}
